package models

type ResponsType MessageType

const (
	Ack ResponsType = "AckResponse"
	Err ResponsType = "ErrorResponse"
)

var responseTypeMap = map[ResponsType]func() CogoMessage{
	Ack: func() CogoMessage { return &AckResponse{BaseCogoMessage: BaseCogoMessage{AckResponseDetails}} },
	Err: func() CogoMessage { return &KillRequest{BaseCogoMessage: BaseCogoMessage{KillRequestDetails}} },
}

func GetResponse(responseType MessageType) CogoMessage {
	if resBuilder := responseTypeMap[ResponsType(responseType)]; resBuilder != nil {
		return resBuilder()
	}
	return nil
}

var AckResponseDetails = CogoMessageDetails{Version: 1, Type: MessageType(Ack)}

type AckResponse struct {
	BaseCogoMessage
}

func NewAckResponse() *AckResponse {
	return &AckResponse{
		BaseCogoMessage: BaseCogoMessage{AckResponseDetails},
	}
}

var ErrResponseDetails = CogoMessageDetails{Version: 1, Type: MessageType(Err)}

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
