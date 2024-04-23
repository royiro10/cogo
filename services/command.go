package services

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/royiro10/cogo/util"
)

type CommandParameters struct {
	SessionId string
	Command   string
}

const DefaultSessionKey = "_"

type CommandService struct {
	logger   *util.Logger
	sessions map[string]*Session
}

func CreateCommandService(logger *util.Logger) *CommandService {
	service := &CommandService{
		logger:   logger,
		sessions: make(map[string]*Session),
	}

	service.sessions[DefaultSessionKey] = NewSession(DefaultSessionKey, service.logger)

	return service
}

func (s *CommandService) HandleCommand(cp *CommandParameters) {
	s.logger.Info("handle command", "sessionId", cp.SessionId, "command", cp.Command)

	session, ok := s.sessions[cp.SessionId]
	if ok {
		s.logger.Debug("valid session Id not provided. using default", "sessionId", cp.SessionId)
		session = s.sessions[DefaultSessionKey]
	}

	args := strings.Fields(cp.Command)

	curser := 0
	for i := 0; i < len(args); i++ {
		if args[i] == "&&" {
			cmd := exec.Command(args[curser], args[curser+1:i]...)
			go session.Run(cmd)
			curser = i + 1
		}
	}

	if curser != len(args) {
		cmd := exec.Command(args[curser], args[curser+1:]...)
		go session.Run(cmd)
	}
}

type Session struct {
	ID string

	stdoutView *bufio.Reader
	stderrView *bufio.Reader
	stdinView  *bufio.Reader

	stdoutChan chan string
	stderrChan chan string
	stdinChan  chan string

	runningCommand *exec.Cmd
	commandQueue   []*exec.Cmd

	logger        *util.Logger
	cancelLogging context.CancelFunc
}

func NewSession(sessionId string, logger *util.Logger) *Session {
	s := &Session{
		ID:           sessionId,
		commandQueue: make([]*exec.Cmd, 0),

		stdoutChan: make(chan string),
		stderrChan: make(chan string),
		stdinChan:  make(chan string),

		logger: logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelLogging = cancel

	if logger != nil {
		s.startStdPipesLogging(logger, ctx)
		logger.Debug("registered pipe logging")
	}

	return s
}

// TODO.com : this is not thread safe. make it use sync.mu
func (s *Session) Run(cmd *exec.Cmd) {
	s.commandQueue = append(s.commandQueue, cmd)
	s.Start()
}

func (s *Session) Start() {
	if s.runningCommand != nil {
		return
	}

	for len(s.commandQueue) > 0 {
		cmd := s.commandQueue[0]
		s.commandQueue = s.commandQueue[1:]
		s.runningCommand = cmd

		if err := s.executeCommand(s.runningCommand); err != nil {
			// TODO: should I handle this???
			s.logger.Error(err.Error(), cmd.String())
			panic(err)
		}

		s.runningCommand = nil
	}
}

type ReplacePipeParameters struct {
	Stdout *bufio.Reader
	Stderr *bufio.Reader
	Stdin  *bufio.Reader
}

func (s *Session) executeCommand(cmd *exec.Cmd) error {
	s.logger.Debug("execution", "command", cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go readPipe(stdout, s.stdoutChan)
	go readPipe(stderr, s.stderrChan)

	err = cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func readPipe(pipe io.ReadCloser, output chan<- string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		if msg := scanner.Text(); msg != "" {
			output <- msg
		}
	}
	pipe.Close()
}

func (s *Session) startStdPipesLogging(logger *util.Logger, ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case out := <-s.stdoutChan:
				logger.Info(out, "reader", "STDOUT")

			case err := <-s.stderrChan:
				logger.Info(err, "reader", "STDERR")

			case in := <-s.stdinChan:
				logger.Info(in, "reader", "STDIN")
			}
		}
	}()
}
