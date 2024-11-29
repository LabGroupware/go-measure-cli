package ws

import (
	"fmt"
)

type CannotParseResponseMessageHandleFunc func(*WebSocket, *CannotParseResponseMessage)

type SubscribeResponseMessageHandleFunc func(*WebSocket, *SubscribeResponseMessage)

type UnsubscribeResponseMessageHandleFunc func(*WebSocket, *UnsubscribeResponseMessage)

type UnsupportedResponseMessageHandleFunc func(*WebSocket, *UnsupportedResponseMessage)

type EventResponseMessageHandleFunc func(*WebSocket, *EventResponseMessage, []byte)

func DefaultCannotParseResponseMessageHandleFunc(ws *WebSocket, msg *CannotParseResponseMessage) {
	fmt.Printf("cannot parse response: %s\n", msg.MessageID)
}

func DefaultSubscribeResponseMessageHandleFunc(ws *WebSocket, msg *SubscribeResponseMessage) {
	fmt.Printf("subscribed: %s\n", msg.SubscriptionID)
}

func DefaultUnsubscribeResponseMessageHandleFunc(ws *WebSocket, msg *UnsubscribeResponseMessage) {
	fmt.Printf("unsubscribed: %s\n", msg.MessageID)
}

func DefaultUnsupportedResponseMessageHandleFunc(ws *WebSocket, msg *UnsupportedResponseMessage) {
	fmt.Printf("unsupported response: %s\n", msg.MessageID)
}

func DefaultEventResponseMessageHandleFunc(ws *WebSocket, msg *EventResponseMessage, raw []byte) {
	fmt.Printf("event response: %s\n", msg.EventType)
}
