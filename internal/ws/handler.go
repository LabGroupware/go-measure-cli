package ws

import (
	"fmt"
)

type CannotParseResponseMessageHandleFunc func(*WebSocket, *CannotParseResponseMessage) error

type SubscribeResponseMessageHandleFunc func(*WebSocket, *SubscribeResponseMessage) error

type UnsubscribeResponseMessageHandleFunc func(*WebSocket, *UnsubscribeResponseMessage) error

type UnsupportedResponseMessageHandleFunc func(*WebSocket, *UnsupportedResponseMessage) error

type EventResponseMessageHandleFunc func(*WebSocket, *EventResponseMessage, []byte) error

func DefaultCannotParseResponseMessageHandleFunc(ws *WebSocket, msg *CannotParseResponseMessage) error {
	fmt.Printf("cannot parse response: %s\n", msg.MessageID)

	return nil
}

func DefaultSubscribeResponseMessageHandleFunc(ws *WebSocket, msg *SubscribeResponseMessage) error {
	fmt.Printf("subscribed: %s\n", msg.SubscriptionID)

	return nil
}

func DefaultUnsubscribeResponseMessageHandleFunc(ws *WebSocket, msg *UnsubscribeResponseMessage) error {
	fmt.Printf("unsubscribed: %s\n", msg.MessageID)

	return nil
}

func DefaultUnsupportedResponseMessageHandleFunc(ws *WebSocket, msg *UnsupportedResponseMessage) error {
	fmt.Printf("unsupported response: %s\n", msg.MessageID)

	return nil
}

func DefaultEventResponseMessageHandleFunc(ws *WebSocket, msg *EventResponseMessage, raw []byte) error {
	fmt.Printf("event response: %s\n", msg.EventType)

	return nil
}
