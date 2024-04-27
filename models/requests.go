package models

type RequestType string

const (
	Execute RequestType = "ExecuteRequest"
	Kill    RequestType = "KillRequest"
)

var requestTypeMap = map[RequestType]func() interface{}{
	Execute: func() interface{} { return &ExecuteRequest{Details: ExecuteRequestDetails} },
	Kill:    func() interface{} { return &KillRequest{Details: KillRequestDetails} },
}

func GetRequest(requestType RequestType) interface{} {
	return requestTypeMap[requestType]()
}

type CogoMessageDetails struct {
	Version int
	Type    RequestType
}

var ExecuteRequestDetails = CogoMessageDetails{Version: 1, Type: Execute}

type ExecuteRequest struct {
	Details   CogoMessageDetails
	SessionId string
	Command   string
}

func NewExecuteRequest(sessionId string, command string) *ExecuteRequest {
	return &ExecuteRequest{
		Details:   ExecuteRequestDetails,
		SessionId: sessionId,
		Command:   command,
	}
}

var KillRequestDetails = CogoMessageDetails{Version: 1, Type: Kill}

type KillRequest struct {
	Details   CogoMessageDetails
	SessionId string
}

func NewKillRequest(sessionId string) *KillRequest {
	return &KillRequest{
		Details:   KillRequestDetails,
		SessionId: sessionId,
	}
}
