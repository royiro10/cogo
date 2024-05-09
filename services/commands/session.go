package commands

import (
	"context"
	"os/exec"
	"sync"
	"time"

	"github.com/royiro10/cogo/common"
	"github.com/royiro10/cogo/models"
)

type Session struct {
	ID string

	queueMu        sync.Mutex
	executionMu    sync.Mutex
	runningCommand *exec.Cmd
	commandQueue   []*exec.Cmd
	killChan       chan struct{}

	stdoutContainer *StdContainer
	stderrContainer *StdContainer
	stdinContainer  *StdContainer

	logger *common.Logger

	cancel context.CancelFunc
}

func NewSession(sessionId string, logger *common.Logger, ctx context.Context) *Session {
	s := &Session{
		ID: sessionId,

		commandQueue: make([]*exec.Cmd, 0),
		killChan:     make(chan struct{}),

		stdoutContainer: NewStdContainer("STDOUT"),
		stderrContainer: NewStdContainer("STDERR"),
		stdinContainer:  NewStdContainer("STDIN"),

		logger: logger,
	}

	sessionCtx, sessionCancelFunc := context.WithCancel(ctx)
	s.cancel = sessionCancelFunc

	s.stdoutContainer.Init(sessionCtx)
	s.stderrContainer.Init(sessionCtx)
	s.stdinContainer.Init(sessionCtx)

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
	go s.startCommandExecution()
}

func (s *Session) Kill() {
	s.stopCommandExecution()
	s.cancel()
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

func (s *Session) startCommandExecution() {
	if !s.executionMu.TryLock() {
		return
	}
	defer s.executionMu.Unlock()

	for s.runningCommand = s.popCommand(); s.runningCommand != nil; s.runningCommand = s.popCommand() {
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
	}
}

func (s *Session) popCommand() *exec.Cmd {
	s.queueMu.Lock()
	defer s.queueMu.Unlock()

	if len(s.commandQueue) < 1 {
		return nil
	}

	cmd := s.commandQueue[0]
	s.commandQueue = s.commandQueue[1:]

	return cmd
}

func (s *Session) stopCommandExecution() {
	s.queueMu.Lock()
	defer s.queueMu.Unlock()

	s.commandQueue = make([]*exec.Cmd, 0)
	if s.runningCommand == nil {
		return
	}

	go func() {
		s.killChan <- struct{}{}
	}()

	runningCommand := s.runningCommand
	s.runningCommand = nil

	if err := runningCommand.Process.Kill(); err != nil {
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
