package commands

import (
	"context"
	"os/exec"
	"strings"
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

	service.sessions[DefaultSessionKey] = NewSession(
		DefaultSessionKey,
		service.logger,
		context.TODO(),
	)

	return service
}

func (s *CommandService) HandleListSessions(request *models.ListSessionsRequest) *models.ListSessionsResponse {
	s.logger.Info("handle list sessions", "sessions", len(s.sessions))
	sessionIds := make([]string, len(s.sessions))

	i := 0
	for k := range s.sessions {
		sessionIds[i] = k
		i++
	}

	return models.NewListSessionsResponse(sessionIds)
}

func (s *CommandService) HandleStatus(request *models.StatusRequest, ctx context.Context) chan *models.SessionStatus {
	s.logger.Info("handle status", "sessionId", request.SessionId, "isStream", request.IsStream)

	statusChan := make(chan *models.SessionStatus)

	session, ok := s.sessions[request.SessionId]
	if !ok {
		s.logger.Warn("no session matching requested session", "session", request.SessionId)
		return statusChan
	}

	switch isStream := request.IsStream; isStream {
	case false:
		go s.getStatusResult(session, statusChan, ctx)
	case true:
		// TODO: implement
		s.logger.Error("not implemented")
	default:
		s.logger.Error("could not recognized output mode")
	}

	return statusChan
}

func (s *CommandService) getStatusResult(session *Session, statusChan chan *models.SessionStatus, ctx context.Context) {
	defer close(statusChan)

	status := &models.SessionStatus{
		SessionStatus:        string(session.GetStatus()),
		LastCommand:          "none",
		LastActionTime:       time.UnixMilli(0),
		ExecutedCommandCount: 0,
		ExecuteQueueSize:     uint(session.GetExecutionQueueSize()),
		OutputViewSize:       uint(session.GetOutputSize()),
	}

	if history := session.GetHistory(); len(history) > 0 {
		status.LastCommand = history[len(history)-1].String()
		status.LastActionTime = session.LastActionTime
		status.ExecutedCommandCount = uint(len(history))
	}

	select {
	case <-ctx.Done():
		s.logger.Info("stop streaming signal was recived")
		return

	default:
		statusChan <- status
	}
}

func (s *CommandService) HandleCommand(request *models.ExecuteRequest) {
	s.logger.Info("Handle command", "sessionId", request.SessionId, "command", request.Command)

	session := s.getOrCreateSession(request.SessionId)

	args := strings.Fields(request.Command)
	commands := make([]*exec.Cmd, 0)

	cursor := 0
	for i := 0; i < len(args); i++ {
		if args[i] == "&&" {
			cmd := exec.Command(args[cursor], args[cursor+1:i]...)
			cmd.Dir = request.Workdir
			commands = append(commands, cmd)
			cursor = i + 1
		}
	}

	if cursor != len(args) {
		cmd := exec.Command(args[cursor], args[cursor+1:]...)
		cmd.Dir = request.Workdir
		commands = append(commands, cmd)
	}

	for _, cmd := range commands {
		session.Run(cmd)
	}
}

func (s *CommandService) HandleKill(request *models.KillRequest) {
	s.logger.Info("Handle kill", "sessionId", request.SessionId)

	session, ok := s.sessions[request.SessionId]
	if !ok {
		s.logger.Warn("no session matching requested session", "session", request.SessionId)
		return
	}

	session.Kill()
}

func (s *CommandService) HandleOutput(
	request *models.OutputRequest,
	ctx context.Context,
) chan models.StdLine {
	s.logger.Info("Handle output", "sessionId", request.SessionId)

	session := s.getOrCreateSession(request.SessionId)
	outputChan := make(chan models.StdLine)

	switch isStream := request.IsStream; isStream {
	case true:
		go s.getOutputStream(session, outputChan, ctx)
	case false:
		go s.getOutputResult(session, outputChan, ctx)
	default:
		s.logger.Error("could not recognized output mode")
	}

	return outputChan
}

func (s *CommandService) getOutputStream(
	session *Session,
	outputChan chan models.StdLine,
	ctx context.Context,
) {
	var notifyStream StdListener = func(line *models.StdLine) {
		outputChan <- *line
	}

	session.stdoutContainer.AddListener(&notifyStream)
	defer session.stdoutContainer.RemoveListener(&notifyStream)

	<-ctx.Done()
	s.logger.Info("Stop streaming signal was recived")
}

func (s *CommandService) getOutputResult(
	session *Session,
	outputChan chan models.StdLine,
	ctx context.Context,
) {
	output := session.GetOutput(-1)
	s.logger.Debug("Output", "view", output)

	defer close(outputChan)

	for lineIndex, line := range *output {
		select {
		case <-ctx.Done():
			s.logger.Info("Stop streaming signal was recived", "outputLineIndex", lineIndex)
			return
		default:
			outputChan <- line
		}
	}
}

func (s *CommandService) getOrCreateSession(sessionId string) *Session {
	session, ok := s.sessions[sessionId]
	if !ok {
		session = NewSession(sessionId, s.logger, context.TODO())
		s.sessions[sessionId] = session

		s.logger.Debug(
			"Requested session Id does not exists. created new session",
			"sessionId",
			sessionId,
		)
	}

	return session
}
