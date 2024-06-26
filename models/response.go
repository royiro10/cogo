package models

type ResponseType MessageType

const (
	AckRes          ResponseType = "AckResponse"
	ErrRes          ResponseType = "ErrorResponse"
	OutputRes       ResponseType = "OutputResponse"
	ListSessionsRes ResponseType = "ListSessionsResponse"
	StatusRes       ResponseType = "StatusResponse"
)

var responseTypeMap = map[ResponseType]func() CogoMessage{
	AckRes: func() CogoMessage {
		return &AckResponse{BaseCogoMessage: BaseCogoMessage{AckResponseDetails}}
	},
	ErrRes: func() CogoMessage {
		return &KillRequest{BaseCogoMessage: BaseCogoMessage{KillRequestDetails}}
	},
	OutputRes: func() CogoMessage {
		return &OutputResponse{BaseCogoMessage: BaseCogoMessage{OutputResponseDetails}}
	},
	ListSessionsRes: func() CogoMessage {
		return &ListSessionsResponse{BaseCogoMessage: BaseCogoMessage{ListSessionsResponseDetails}}
	},
	StatusRes: func() CogoMessage {
		return &StatusResponse{BaseCogoMessage: BaseCogoMessage{StatusResponseDetails}}
	},
}

func GetResponse(responseType MessageType) CogoMessage {
	if resBuilder := responseTypeMap[ResponseType(responseType)]; resBuilder != nil {
		return resBuilder()
	}
	return nil
}

var AckResponseDetails = CogoMessageDetails{Version: 1, Type: MessageType(AckRes)}

type AckResponse struct {
	BaseCogoMessage
}

func NewAckResponse() *AckResponse {
	return &AckResponse{
		BaseCogoMessage: BaseCogoMessage{AckResponseDetails},
	}
}

var ErrResponseDetails = CogoMessageDetails{Version: 1, Type: MessageType(ErrRes)}

type ErrResponse struct {
	BaseCogoMessage
	Err error
}

func NewErrResponse(err error) *ErrResponse {
	return &ErrResponse{
		BaseCogoMessage: BaseCogoMessage{ExecuteRequestDetails},
		Err:             err,
	}
}

var OutputResponseDetails = CogoMessageDetails{Version: 1, Type: MessageType(OutputRes)}

type OutputResponse struct {
	BaseCogoMessage
	StdLine
}

func NewOutputResponse(output *StdLine) *OutputResponse {
	return &OutputResponse{
		BaseCogoMessage: BaseCogoMessage{OutputResponseDetails},
		StdLine:         *output,
	}
}

var ListSessionsResponseDetails = CogoMessageDetails{Version: 1, Type: MessageType(ListSessionsRes)}

type ListSessionsResponse struct {
	BaseCogoMessage
	Sessions []string
}

func NewListSessionsResponse(sessions []string) *ListSessionsResponse {
	return &ListSessionsResponse{
		BaseCogoMessage: BaseCogoMessage{ListSessionsResponseDetails},
		Sessions:        sessions,
	}
}

var StatusResponseDetails = CogoMessageDetails{Version: 1, Type: MessageType(StatusRes)}

type StatusResponse struct {
	BaseCogoMessage
	SessionId string
	Status    SessionStatus
}

func NewStatusResponse(sessionId string, status SessionStatus) *StatusResponse {
	return &StatusResponse{
		BaseCogoMessage: BaseCogoMessage{StatusResponseDetails},
		SessionId:       sessionId,
		Status:          status,
	}
}
