package ws

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/LabGroupware/go-measure-tui/internal/utils"
	"github.com/gorilla/websocket"
)

// WebSocket is a struct that represents a WebSocket connection
type WebSocket struct {
	conn                  *websocket.Conn
	CannotParseMsgHandler CannotParseResponseMessageHandleFunc
	SubscribeMsgHandler   SubscribeResponseMessageHandleFunc
	UnsubscribeMsgHandler UnsubscribeResponseMessageHandleFunc
	UnsupportedMsgHandler UnsupportedResponseMessageHandleFunc
	EventMsgHandler       EventResponseMessageHandleFunc
	subscribeMsgMem       *SubscribeMessageMemory
	unsubscribeMsgMem     *UnsubscribeMessageMemory
	SubscribeMemory       *SubscribeMemory
}

// NewWebSocket creates a new WebSocket connection
func NewWebSocket() *WebSocket {
	return &WebSocket{
		conn:                  nil,
		CannotParseMsgHandler: DefaultCannotParseResponseMessageHandleFunc,
		SubscribeMsgHandler:   DefaultSubscribeResponseMessageHandleFunc,
		UnsubscribeMsgHandler: DefaultUnsubscribeResponseMessageHandleFunc,
		UnsupportedMsgHandler: DefaultUnsupportedResponseMessageHandleFunc,
		EventMsgHandler:       DefaultEventResponseMessageHandleFunc,
	}
}

type ConnectConfig struct {
	DisconnectOnReadMsgError       bool
	DisconnectOnUnmarshalJSONError bool
	DisconnectOnParseMsgError      bool
}

// Connect connects to a remote server using the WebSocket protocol
func (ws *WebSocket) Connect(ctx context.Context, ctr *app.Container, conf ConnectConfig) (<-chan struct{}, error) {
	headers := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", ctr.AuthToken.AccessToken)},
	}
	conn, res, err := websocket.DefaultDialer.Dial(ctr.Config.Web.WebSocket.Url, headers)
	if err != nil {
		if res != nil && res.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("unauthorized: %w", err)
		}
		return nil, fmt.Errorf("failed to connect to WebSocket server: %w", err)
	}
	ws.conn = conn

	// Clean up the connection when the function returns
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	go func() {
		defer close(done)
		defer ws.conn.Close()
		for {
			var msg ResponseMessage
			_, content, err := ws.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					ctr.Logger.Debug(ctx, "WebSocket connection closed",
						logger.Value("error", err))
					return
				}
				ctr.Logger.Error(ctx, "failed to read message",
					logger.Value("error", err))
				if conf.DisconnectOnReadMsgError {
					return
				}
			}
			err = utils.UnmarshalJSON(content, &msg)
			if err != nil {
				ctr.Logger.Error(ctx, "failed to unmarshal JSON",
					logger.Value("error", err))
				if conf.DisconnectOnUnmarshalJSONError {
					return
				}
			}
			switch msg.Type {
			case ResponseTypeCannotParseResponse:
				var cannotParse CannotParseResponseMessage
				err := utils.UnmarshalJSON(content, &cannotParse)
				if err != nil {
					ctr.Logger.Error(ctx, "failed to read CannotParseResponseMessage",
						logger.Value("error", err))
					if conf.DisconnectOnParseMsgError {
						return
					}
				}
				if ws.CannotParseMsgHandler != nil {
					ws.CannotParseMsgHandler(ws, &cannotParse)
				}
			case ResponseTypeSubscribeResponse:
				var subscribe SubscribeResponseMessage
				err := utils.UnmarshalJSON(content, &subscribe)
				if err != nil {
					ctr.Logger.Error(ctx, "failed to read SubscribeResponseMessage",
						logger.Value("error", err))
					if conf.DisconnectOnParseMsgError {
						return
					}
				}
				if v, ok := ws.subscribeMsgMem.Memory[subscribe.MessageID]; ok {
					if subscribe.Success {
						ws.SubscribeMemory.Memory[subscribe.SubscriptionID] = SubscribeMemoryData(v)
						delete(ws.subscribeMsgMem.Memory, subscribe.MessageID)
					}
				}
				if ws.SubscribeMsgHandler != nil {
					ws.SubscribeMsgHandler(ws, &subscribe)
				}
			case ResponseTypeUnsubscribeResponse:
				var unsubscribe UnsubscribeResponseMessage
				err := utils.UnmarshalJSON(content, &unsubscribe)
				if err != nil {
					ctr.Logger.Error(ctx, "failed to read UnsubscribeResponseMessage",
						logger.Value("error", err))
					if conf.DisconnectOnParseMsgError {
						return
					}
				}
				if v, ok := ws.unsubscribeMsgMem.Memory[unsubscribe.MessageID]; ok {
					if unsubscribe.Success {
						for _, id := range v.SubscriptionIDs {
							delete(ws.SubscribeMemory.Memory, id)
						}
						delete(ws.unsubscribeMsgMem.Memory, unsubscribe.MessageID)
					}
				}
				if ws.UnsubscribeMsgHandler != nil {
					ws.UnsubscribeMsgHandler(ws, &unsubscribe)
				}
			case ResponseTypeUnsupportedResponse:
				var unsupported UnsupportedResponseMessage
				err := utils.UnmarshalJSON(content, &unsupported)
				if err != nil {
					ctr.Logger.Error(ctx, "failed to read UnsupportedResponseMessage",
						logger.Value("error", err))
					if conf.DisconnectOnParseMsgError {
						return
					}
				}
				if ws.UnsupportedMsgHandler != nil {
					ws.UnsupportedMsgHandler(ws, &unsupported)
				}
			case ResponseTypeEventResponse:
				var event EventResponseMessage
				err := utils.UnmarshalJSON(content, &event)
				if err != nil {
					ctr.Logger.Error(ctx, "failed to read EventResponseMessage",
						logger.Value("error", err))
					if conf.DisconnectOnParseMsgError {
						return
					}
				}
				if ws.EventMsgHandler != nil {
					ws.EventMsgHandler(ws, &event, content)
				}
			}
		}
	}()

	return done, nil
}

func (ws *WebSocket) AllUnsubscribe() error {
	var subscriptionIDs []string
	for id := range ws.SubscribeMemory.Memory {
		subscriptionIDs = append(subscriptionIDs, id)
	}
	if len(subscriptionIDs) == 0 {
		return nil
	}
	return ws.SendUnsubscribeMessage(subscriptionIDs)
}

func (ws *WebSocket) SendSubscribeMessage(aggregateType AggregateType, aggregateIDs []string, eventTypes []EventType) error {
	msgID := utils.GenerateUniqueID()
	msg := SubscribeMessage{
		Type:      RequestTypeSubscribe,
		MessageID: msgID,
		Data:      SubscribeData{AggregateType: aggregateType, AggregateIDs: aggregateIDs, EventTypes: eventTypes},
	}
	err := ws.conn.WriteJSON(msg)
	if err != nil {
		return fmt.Errorf("failed to send subscribe message: %w", err)
	}
	ws.subscribeMsgMem.Memory[msgID] = SubscribeMessageMemoryData{AggregateType: aggregateType, AggregateIDs: aggregateIDs, EventTypes: eventTypes}
	return nil
}

func (ws *WebSocket) SendUnsubscribeMessage(subscriptionIDs []string) error {
	msgID := utils.GenerateUniqueID()
	msg := UnsubscribeMessage{
		Type:      RequestTypeUnsubscribe,
		MessageID: msgID,
		Data:      UnsubscribeData{SubscriptionIDs: subscriptionIDs},
	}
	err := ws.conn.WriteJSON(msg)
	if err != nil {
		return fmt.Errorf("failed to send unsubscribe message: %w", err)
	}
	ws.unsubscribeMsgMem.Memory[msgID] = UnsubscribeMessageMemoryData{SubscriptionIDs: subscriptionIDs}
	return nil
}

func (ws *WebSocket) UnsubscribeFromAnyAggregate(aggregateType AggregateType, aggregateID []string) error {
	var subscriptionIDs []string
	for id, v := range ws.SubscribeMemory.Memory {
		if v.AggregateType == aggregateType && utils.AnyContains(v.AggregateIDs, aggregateID) {
			subscriptionIDs = append(subscriptionIDs, id)
		}
	}
	if len(subscriptionIDs) == 0 {
		return nil
	}
	return ws.SendUnsubscribeMessage(subscriptionIDs)
}

func (ws *WebSocket) UnsubscribeFromAllAggregate(aggregateType AggregateType, aggregateID []string) error {
	var subscriptionIDs []string
	for id, v := range ws.SubscribeMemory.Memory {
		if v.AggregateType == aggregateType && utils.AllContains(v.AggregateIDs, aggregateID) {
			subscriptionIDs = append(subscriptionIDs, id)
		}
	}
	if len(subscriptionIDs) == 0 {
		return nil
	}
	return ws.SendUnsubscribeMessage(subscriptionIDs)
}

func (ws *WebSocket) UnsubscribeFromAnyEvent(eventType []EventType) error {
	var subscriptionIDs []string
	for id, v := range ws.SubscribeMemory.Memory {
		if utils.AnyContains(v.EventTypes, eventType) {
			subscriptionIDs = append(subscriptionIDs, id)
		}
	}
	if len(subscriptionIDs) == 0 {
		return nil
	}
	return ws.SendUnsubscribeMessage(subscriptionIDs)
}

func (ws *WebSocket) UnsubscribeFromAllEvent(eventType []EventType) error {
	var subscriptionIDs []string
	for id, v := range ws.SubscribeMemory.Memory {
		if utils.AllContains(v.EventTypes, eventType) {
			subscriptionIDs = append(subscriptionIDs, id)
		}
	}
	if len(subscriptionIDs) == 0 {
		return nil
	}
	return ws.SendUnsubscribeMessage(subscriptionIDs)
}

func (ws *WebSocket) UnsubscribeFromAnyAggregateAndAnyEvent(aggregateType AggregateType, aggregateID []string, eventType []EventType) error {
	var subscriptionIDs []string
	for id, v := range ws.SubscribeMemory.Memory {
		if v.AggregateType == aggregateType && utils.AnyContains(v.AggregateIDs, aggregateID) && utils.AnyContains(v.EventTypes, eventType) {
			subscriptionIDs = append(subscriptionIDs, id)
		}
	}
	if len(subscriptionIDs) == 0 {
		return nil
	}
	return ws.SendUnsubscribeMessage(subscriptionIDs)
}

func (ws *WebSocket) UnsubscribeFromAnyAggregateAndAllEvent(aggregateType AggregateType, aggregateID []string, eventType []EventType) error {
	var subscriptionIDs []string
	for id, v := range ws.SubscribeMemory.Memory {
		if v.AggregateType == aggregateType && utils.AnyContains(v.AggregateIDs, aggregateID) && utils.AllContains(v.EventTypes, eventType) {
			subscriptionIDs = append(subscriptionIDs, id)
		}
	}
	if len(subscriptionIDs) == 0 {
		return nil
	}
	return ws.SendUnsubscribeMessage(subscriptionIDs)
}

func (ws *WebSocket) UnsubscribeFromAllAggregateAndAnyEvent(aggregateType AggregateType, aggregateID []string, eventType []EventType) error {
	var subscriptionIDs []string
	for id, v := range ws.SubscribeMemory.Memory {
		if v.AggregateType == aggregateType && utils.AllContains(v.AggregateIDs, aggregateID) && utils.AnyContains(v.EventTypes, eventType) {
			subscriptionIDs = append(subscriptionIDs, id)
		}
	}
	if len(subscriptionIDs) == 0 {
		return nil
	}
	return ws.SendUnsubscribeMessage(subscriptionIDs)
}

func (ws *WebSocket) UnsubscribeFromAllAggregateAndAllEvent(aggregateType AggregateType, aggregateID []string, eventType []EventType) error {
	var subscriptionIDs []string
	for id, v := range ws.SubscribeMemory.Memory {
		if v.AggregateType == aggregateType && utils.AllContains(v.AggregateIDs, aggregateID) && utils.AllContains(v.EventTypes, eventType) {
			subscriptionIDs = append(subscriptionIDs, id)
		}
	}
	if len(subscriptionIDs) == 0 {
		return nil
	}
	return ws.SendUnsubscribeMessage(subscriptionIDs)
}

// Send CloseMessage sends a close message to the server
func (ws *WebSocket) SendCloseMessage() error {
	err := ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return fmt.Errorf("failed to send close message: %w", err)
	}
	return nil
}

// Close closes the WebSocket connection
func (ws *WebSocket) Close() {
	ws.conn.Close()
}
