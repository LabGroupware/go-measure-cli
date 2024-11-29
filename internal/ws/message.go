package ws

type Message struct {
	MessageID string      `json:"messageId"`
	Type      RequestType `json:"type"`
	Data      any         `json:"data"`
}

type SendMessage[T any] struct {
	MessageID string      `json:"messageId"`
	Type      RequestType `json:"type"`
	Data      T           `json:"data"`
}

type SubscribeData struct {
	AggregateType AggregateType `json:"aggregateType"`
	AggregateIDs  []string      `json:"aggregateIds"`
	EventTypes    []EventType   `json:"eventTypes"`
}

type SubscribeMessage SendMessage[SubscribeData]

type UnsubscribeData struct {
	SubscriptionIDs []string `json:"subscriptionIds"`
}

type UnsubscribeMessage SendMessage[UnsubscribeData]
