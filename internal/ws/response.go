package ws

type ResponseMessage struct {
	Type ResponseType `json:"type"`
}

type CannotParseResponseMessage struct {
	MessageID   string       `json:"messageId"`
	Type        ResponseType `json:"type"`
	RequestType string       `json:"requestType"`
	Success     bool         `json:"success"`
}

type SubscribeResponseMessage struct {
	MessageID      string       `json:"messageId"`
	Type           ResponseType `json:"type"`
	SubscriptionID string       `json:"subscriptionId"`
	Success        bool         `json:"success"`
}

type UnsubscribeResponseMessage struct {
	MessageID string       `json:"messageId"`
	Type      ResponseType `json:"type"`
	Success   bool         `json:"success"`
}

type UnsupportedResponseMessage struct {
	MessageID string       `json:"messageId"`
	Type      ResponseType `json:"type"`
	Success   bool         `json:"success"`
}

type EventResponseMessage struct {
	Type      ResponseType `json:"type"`
	Data      any          `json:"data"`
	EventType EventType    `json:"eventType"`
}

type EventResponseMessageWithData[T any] struct {
	Type      ResponseType `json:"type"`
	Data      T            `json:"data"`
	EventType EventType    `json:"eventType"`
}
