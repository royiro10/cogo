package models

type RequestType MessageType

const (
	ExecuteReq      RequestType = "ExecuteRequest"
	KillReq         RequestType = "KillRequest"
	OutputReq       RequestType = "OutputRequest"
	StatusReq       RequestType = "StatusRequest"
	ListSessionsReq RequestType = "ListSessionsRequest"
)

type RequestBuilder func() CogoMessage

var requestTypeMap = map[RequestType]RequestBuilder{
	ExecuteReq: func() CogoMessage {
		return &ExecuteRequest{BaseCogoMessage: BaseCogoMessage{ExecuteRequestDetails}}
	},
	KillReq: func() CogoMessage {
		return &KillRequest{BaseCogoMessage: BaseCogoMessage{KillRequestDetails}}
	},
	OutputReq: func() CogoMessage {
		return &OutputRequest{BaseCogoMessage: BaseCogoMessage{OutputRequestDetails}}
	},
	StatusReq: func() CogoMessage {
		return &StatusRequest{BaseCogoMessage: BaseCogoMessage{StatusRequestDetails}}
	},
	ListSessionsReq: func() CogoMessage {
		return &ListSessionsRequest{BaseCogoMessage: BaseCogoMessage{ListSessionRequestDeatils}}
	},
}

func GetRequest(requestType MessageType) CogoMessage {
	if reqBuilder := requestTypeMap[RequestType(requestType)]; reqBuilder != nil {
		return reqBuilder()
	}
	return nil
}

var ExecuteRequestDetails = CogoMessageDetails{Version: 1, Type: MessageType(ExecuteReq)}

type ExecuteRequest struct {
	BaseCogoMessage
	SessionId string
	Command   string
	Workdir   string
	IsRestart bool
}

func NewExecuteRequest(
	sessionId string,
	command string,
	workdir string,
	isRestart bool,
) *ExecuteRequest {
	return &ExecuteRequest{
		BaseCogoMessage: BaseCogoMessage{ExecuteRequestDetails},
		SessionId:       sessionId,
		Command:         command,
		Workdir:         workdir,
		IsRestart:       isRestart,
	}
}

var KillRequestDetails = CogoMessageDetails{Version: 1, Type: MessageType(KillReq)}

type KillRequest struct {
	BaseCogoMessage
	SessionId string
}

func NewKillRequest(sessionId string) *KillRequest {
	return &KillRequest{
		BaseCogoMessage: BaseCogoMessage{KillRequestDetails},
		SessionId:       sessionId,
	}
}

var OutputRequestDetails = CogoMessageDetails{Version: 1, Type: MessageType(OutputReq)}

type OutputRequest struct {
	BaseCogoMessage
	SessionId string
	IsStream  bool
}

func NewOutputRequest(sessionId string, isStream bool) *OutputRequest {
	return &OutputRequest{
		BaseCogoMessage: BaseCogoMessage{OutputRequestDetails},
		SessionId:       sessionId,
		IsStream:        isStream,
	}
}

var ListSessionRequestDeatils = CogoMessageDetails{Version: 1, Type: MessageType(ListSessionsReq)}

type ListSessionsRequest struct {
	BaseCogoMessage
}

func NewListSessionsRequest() *ListSessionsRequest {
	return &ListSessionsRequest{
		BaseCogoMessage: BaseCogoMessage{ListSessionRequestDeatils},
	}
}

var StatusRequestDetails = CogoMessageDetails{Version: 1, Type: MessageType(StatusReq)}

type StatusRequest struct {
	BaseCogoMessage
	SessionId string
	IsStream  bool
}

func NewStatusRequest(sessionId string, isStream bool) *StatusRequest {
	return &StatusRequest{
		BaseCogoMessage: BaseCogoMessage{StatusRequestDetails},
		SessionId:       sessionId,
		IsStream:        isStream,
	}
}
