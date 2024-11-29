package ws

type RequestType string

const (
	RequestTypeSubscribe   RequestType = "subscribe"
	RequestTypeUnsubscribe RequestType = "unsubscribe"
)
