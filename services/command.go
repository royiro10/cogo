package services

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

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
		session.Run(cmd)
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

func (s *CommandService) HandleOutput(request *models.OutputRequest) *[]models.StdLine {
	s.logger.Info("handle output", "sessionId", request.SessionId)

	session, ok := s.sessions[request.SessionId]
	if !ok {
		session = NewSession(request.SessionId, s.logger)
		s.sessions[request.SessionId] = session

		s.logger.Debug("valid session Id not provided. created new session", "sessionId", request.SessionId)
	}

	output := session.GetOutput(-1)
	s.logger.Info("output", "view", output)
	return output
}

type Session struct {
	ID       string
	ExecChan chan *exec.Cmd

	queueMu        sync.Mutex
	executionMu    sync.Mutex
	runningCommand *exec.Cmd
	commandQueue   []*exec.Cmd
	killChan       chan struct{}

	stdoutContainer *StdContainer
	stderrContainer *StdContainer
	stdinContainer  *StdContainer

	logger        *common.Logger
	cancelLogging context.CancelFunc
}

func NewSession(sessionId string, logger *common.Logger) *Session {
	s := &Session{
		ID:       sessionId,
		ExecChan: make(chan *exec.Cmd),

		commandQueue: make([]*exec.Cmd, 0),
		killChan:     make(chan struct{}),

		stdoutContainer: NewStdContainer("STDOUT"),
		stderrContainer: NewStdContainer("STDERR"),
		stdinContainer:  NewStdContainer("STDIN"),

		logger: logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelLogging = cancel

	s.stdoutContainer.Init(ctx)
	s.stderrContainer.Init(ctx)
	s.stdinContainer.Init(ctx)

	if logger != nil {
		s.stdoutContainer.AddListener(makePipeLogger(s.stdoutContainer, logger))
		s.stderrContainer.AddListener(makePipeLogger(s.stderrContainer, logger))
		s.stdinContainer.AddListener(makePipeLogger(s.stdinContainer, logger))

		logger.Debug("registered pipe logging")
	}

	return s
}

func (s *Session) Run(cmd *exec.Cmd) {
	s.queueMu.Lock()
	defer s.queueMu.Unlock()

	s.commandQueue = append(s.commandQueue, cmd)
	go s.Start()
}

func (s *Session) Kill() {
	s.Stop()
	s.cancelLogging()
	s.logger.Info("killed session", "sessionId", s.ID)
}

func (s *Session) GetOutput(tailCount int) *[]models.StdLine {
	if tailCount == -1 {
		output := s.stdoutContainer.View()
		return &output
	}

	output := s.stdoutContainer.ViewTail(tailCount)
	return &output
}

func (s *Session) Start() {
	if !s.executionMu.TryLock() {
		return
	}
	defer s.executionMu.Unlock()

	for len(s.commandQueue) > 0 {
		s.queueMu.Lock()
		cmd := s.commandQueue[0]
		s.commandQueue = s.commandQueue[1:]
		s.runningCommand = cmd
		s.queueMu.Unlock()

		err := s.executeCommand(s.runningCommand)
		if err != nil {
			select {
			case <-s.killChan:

			default:
				s.stderrContainer.NotifyChan <- models.StdLine{
					Time: time.Now(),
					Data: err.Error(),
				}
			}
		}

		s.runningCommand = nil
	}
}

func (s *Session) Stop() {
	s.commandQueue = make([]*exec.Cmd, 0)

	go func() {
		s.killChan <- struct{}{}
	}()

	if err := s.runningCommand.Process.Kill(); err != nil {
		s.logger.Warn("error while canceling command", "err", err)
		return
	}
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

	go readPipe(stdout, s.stdoutContainer.NotifyChan)
	go readPipe(stderr, s.stderrContainer.NotifyChan)

	err = cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func readPipe(pipe io.ReadCloser, output chan<- models.StdLine) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		if msg := scanner.Text(); msg != "" {
			output <- models.StdLine{
				Time: time.Now(),
				Data: msg,
			}
		}
	}

	pipe.Close()
}

func makePipeLogger(sc *StdContainer, logger *common.Logger) *StdListener {
	var pipeLoggerListener StdListener = func(line *models.StdLine) {
		logger.Info(line.Data, "reader", sc.Name)
	}

	return &pipeLoggerListener
}
