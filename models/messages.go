package models

type MessageType string

type CogoMessageDetails struct {
	Version int
	Type    MessageType
}

type CogoMessage interface {
	GetDetails() *CogoMessageDetails
}

type BaseCogoMessage struct {
	Details CogoMessageDetails
}

func (msg *BaseCogoMessage) GetDetails() *CogoMessageDetails {
	return &msg.Details
}

func GetMessage(details CogoMessageDetails) CogoMessage {
	if req := GetRequest(details.Type); req != nil {
		return req
	}

	if res := GetResponse(details.Type); res != nil {
		return res
	}

	return nil
}
