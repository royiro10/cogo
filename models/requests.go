package models

type RequestType MessageType

const (
	ExecuteReq RequestType = "ExecuteRequest"
	KillReq    RequestType = "KillRequest"
	OutputReq  RequestType = "OutputRequest"
)

type RequestBuilder func() CogoMessage

var requestTypeMap = map[RequestType]RequestBuilder{
	ExecuteReq: func() CogoMessage { return &ExecuteRequest{BaseCogoMessage: BaseCogoMessage{ExecuteRequestDetails}} },
	KillReq:    func() CogoMessage { return &KillRequest{BaseCogoMessage: BaseCogoMessage{KillRequestDetails}} },
	OutputReq:  func() CogoMessage { return &OutputRequest{BaseCogoMessage: BaseCogoMessage{OutputRequestDetails}} },
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
}

func NewExecuteRequest(sessionId string, command string) *ExecuteRequest {
	return &ExecuteRequest{
		BaseCogoMessage: BaseCogoMessage{ExecuteRequestDetails},
		SessionId:       sessionId,
		Command:         command,
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
