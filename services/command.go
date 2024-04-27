package services

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
)

const DefaultSessionKey = "_"

type CommandService struct {
	logger   *common.Logger
	sessions map[string]*Session
}

func CreateCommandService(logger *common.Logger) *CommandService {
	service := &CommandService{
		logger:   logger,
		sessions: make(map[string]*Session),
	}

	service.sessions[DefaultSessionKey] = NewSession(DefaultSessionKey, service.logger)

	return service
}

func (s *CommandService) HandleCommand(request *models.ExecuteRequest) {
	s.logger.Info("handle command", "sessionId", request.SessionId, "command", request.Command)

	session, ok := s.sessions[request.SessionId]
	if !ok {
		session = NewSession(request.SessionId, s.logger)
		s.sessions[request.SessionId] = session

		s.logger.Debug("valid session Id not provided. created new session", "sessionId", request.SessionId)
	}

	args := strings.Fields(request.Command)

	commands := make([]*exec.Cmd, 0)

	curser := 0
	for i := 0; i < len(args); i++ {
		if args[i] == "&&" {
			commands = append(commands, exec.Command(args[curser], args[curser+1:i]...))
			curser = i + 1
		}
	}

	if curser != len(args) {
		commands = append(commands, exec.Command(args[curser], args[curser+1:]...))
	}

	for _, cmd := range commands {
		go session.Run(cmd)
	}
}

func (s *CommandService) HandleKill(request *models.KillRequest) {
	s.logger.Info("handle kill", "sessionId", request.SessionId)

	session, ok := s.sessions[request.SessionId]
	if !ok {
		s.logger.Warn("no session matching requested session", "session", request.SessionId)
		return
	}

	session.Kill()
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
	killChan       chan struct{}

	logger        *common.Logger
	cancelLogging context.CancelFunc
}

func NewSession(sessionId string, logger *common.Logger) *Session {
	s := &Session{
		ID:           sessionId,
		commandQueue: make([]*exec.Cmd, 0),

		stdoutChan: make(chan string),
		stderrChan: make(chan string),
		stdinChan:  make(chan string),

		killChan: make(chan struct{}),

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

// TODO : this is not thread safe. make it use sync.mu
func (s *Session) Run(cmd *exec.Cmd) {
	s.commandQueue = append(s.commandQueue, cmd)
	s.Start()
}

func (s *Session) Kill() {
	s.Stop()
	s.cancelLogging()
	s.logger.Info("killed session", "sessionId", s.ID)
}

// TODO : this is not thread safe. make it use sync.mu
func (s *Session) Start() {
	if s.runningCommand != nil {
		return
	}

	for len(s.commandQueue) > 0 {
		cmd := s.commandQueue[0]
		s.commandQueue = s.commandQueue[1:]
		s.runningCommand = cmd

		// TODO: should I handle this???
		err := s.executeCommand(s.runningCommand)
		if err != nil {
			select {
			case <-s.killChan:

			default:
				s.logger.Fatal(common.WrapedError{Msg: cmd.String(), Err: err})
			}
		}

		s.runningCommand = nil
	}
}

func (s *Session) Stop() {
	s.commandQueue = make([]*exec.Cmd, 0)

	s.killChan <- struct{}{}
	if err := s.runningCommand.Process.Kill(); err != nil {
		s.logger.Warn("error while canceling command", "err", err)
		return
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

func (s *Session) startStdPipesLogging(logger *common.Logger, ctx context.Context) {
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
