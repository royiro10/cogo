package models

type ResponseType MessageType

const (
	AckRes    ResponseType = "AckResponse"
	ErrRes    ResponseType = "ErrorResponse"
	OutputRes ResponseType = "OutputResponse"
)

var responseTypeMap = map[ResponseType]func() CogoMessage{
	AckRes:    func() CogoMessage { return &AckResponse{BaseCogoMessage: BaseCogoMessage{AckResponseDetails}} },
	ErrRes:    func() CogoMessage { return &KillRequest{BaseCogoMessage: BaseCogoMessage{KillRequestDetails}} },
	OutputRes: func() CogoMessage { return &OutputResponse{BaseCogoMessage: BaseCogoMessage{OutputResponseDetails}} },
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
	Lines []StdLine
}

func NewOutputResponse(output *[]StdLine) *OutputResponse {
	return &OutputResponse{
		BaseCogoMessage: BaseCogoMessage{OutputResponseDetails},
		Lines:           *output,
	}
}
