package models

type RequestType MessageType

const (
	Execute RequestType = "ExecuteRequest"
	Kill    RequestType = "KillRequest"
)

var requestTypeMap = map[RequestType]func() CogoMessage{
	Execute: func() CogoMessage { return &ExecuteRequest{BaseCogoMessage: BaseCogoMessage{ExecuteRequestDetails}} },
	Kill:    func() CogoMessage { return &KillRequest{BaseCogoMessage: BaseCogoMessage{KillRequestDetails}} },
}

func GetRequest(requestType MessageType) CogoMessage {
	if reqBuilder := requestTypeMap[RequestType(requestType)]; reqBuilder != nil {
		return reqBuilder()
	}
	return nil
}

var ExecuteRequestDetails = CogoMessageDetails{Version: 1, Type: MessageType(Execute)}

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

var KillRequestDetails = CogoMessageDetails{Version: 1, Type: MessageType(Kill)}

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
