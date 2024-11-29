package ws

type ResponseType string

const (
	ResponseTypeCannotParseResponse ResponseType = "cannot-parse"
	ResponseTypeSubscribeResponse   ResponseType = "subscribe"
	ResponseTypeUnsubscribeResponse ResponseType = "unsubscribe"
	ResponseTypeUnsupportedResponse ResponseType = "unsupported"
	ResponseTypeEventResponse       ResponseType = "event"
)
