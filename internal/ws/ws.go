package ws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/LabGroupware/go-measure-tui/internal/utils"
	"github.com/gorilla/websocket"
)

// WebSocket is a struct that represents a WebSocket connection
type WebSocket struct {
	conn                            *websocket.Conn
	CannotParseMsgHandler           CannotParseResponseMessageHandleFunc
	SubscribeMsgHandler             SubscribeResponseMessageHandleFunc
	UnsubscribeMsgHandler           UnsubscribeResponseMessageHandleFunc
	UnsupportedMsgHandler           UnsupportedResponseMessageHandleFunc
	EventMsgHandler                 EventResponseMessageHandleFunc
	subscribeMsgMem                 *SubscribeMessageMemory
	unsubscribeMsgMem               *UnsubscribeMessageMemory
	SubscribeMemory                 *SubscribeMemory
	SubscribeConsumerMemory         *SubscribeConsumerMemory
	SubscribeConsumerMemorySubIndex *SubscribeConsumerMemorySubIndex
}

// NewWebSocket creates a new WebSocket connection
func NewWebSocket() *WebSocket {
	return &WebSocket{
		conn:                            nil,
		CannotParseMsgHandler:           DefaultCannotParseResponseMessageHandleFunc,
		SubscribeMsgHandler:             DefaultSubscribeResponseMessageHandleFunc,
		UnsubscribeMsgHandler:           DefaultUnsubscribeResponseMessageHandleFunc,
		UnsupportedMsgHandler:           DefaultUnsupportedResponseMessageHandleFunc,
		EventMsgHandler:                 DefaultEventResponseMessageHandleFunc,
		subscribeMsgMem:                 NewSubscribeMessageMemory(),
		unsubscribeMsgMem:               NewUnsubscribeMessageMemory(),
		SubscribeMemory:                 NewSubscribeMemory(),
		SubscribeConsumerMemory:         NewSubscribeConsumerMemory(),
		SubscribeConsumerMemorySubIndex: NewSubscribeConsumerMemorySubIndex(),
	}
}

type ConnectConfig struct {
	DisconnectOnReadMsgError       bool
	DisconnectOnUnmarshalJSONError bool
	DisconnectOnParseMsgError      bool
}

// Connect connects to a remote server using the WebSocket protocol
func (ws *WebSocket) Connect(ctx context.Context, ctr *app.Container, conf ConnectConfig) (<-chan TerminateType, error) {
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
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)

	done := make(chan TerminateType)

	go func() {
		// defer close(done)
		// defer ws.conn.Close()
		for {
			var msg ResponseMessage
			_, content, err := ws.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					ctr.Logger.Debug(ctx, "WebSocket connection closed",
						logger.Value("error", err))
					done <- TerminateTypeConnectionClosed
					return
				}
				ctr.Logger.Error(ctx, "failed to read message",
					logger.Value("error", err))
				if conf.DisconnectOnReadMsgError {
					done <- TerminateTypeReadMsgError
					return
				}
			}
			err = utils.UnmarshalJSON(content, &msg)
			if err != nil {
				ctr.Logger.Error(ctx, "failed to unmarshal JSON",
					logger.Value("error", err))
				if conf.DisconnectOnUnmarshalJSONError {
					done <- TerminateTypeUnmarshalJSONError
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
						done <- TerminateTypeParseMsgError
						return
					}
				}
				if ws.CannotParseMsgHandler != nil {
					if err := ws.CannotParseMsgHandler(ws, &cannotParse); err != nil {
						ctr.Logger.Error(ctx, "failed to handle CannotParseResponseMessage",
							logger.Value("error", err))
						done <- TerminateTypeCannotParseMsgHandlerError
						return
					}
				}
			case ResponseTypeSubscribeResponse:
				var subscribe SubscribeResponseMessage
				err := utils.UnmarshalJSON(content, &subscribe)
				if err != nil {
					ctr.Logger.Error(ctx, "failed to read SubscribeResponseMessage",
						logger.Value("error", err))
					if conf.DisconnectOnParseMsgError {
						done <- TerminateTypeParseMsgError
						return
					}
				}
				if v, ok := ws.subscribeMsgMem.Memory[subscribe.MessageID]; ok {
					if subscribe.Success {
						ws.SubscribeMemory.Memory[subscribe.SubscriptionID] = SubscribeMemoryData{
							AggregateType: v.AggregateType,
							AggregateIDs:  v.AggregateIDs,
							EventTypes:    v.EventTypes,
						}
						ws.SubscribeConsumerMemory.Memory[v.ConsumerID] = subscribe.SubscriptionID
						ws.SubscribeConsumerMemorySubIndex.Memory[subscribe.SubscriptionID] = v.ConsumerID
						delete(ws.subscribeMsgMem.Memory, subscribe.MessageID)
					}
				}
				if ws.SubscribeMsgHandler != nil {
					if err := ws.SubscribeMsgHandler(ws, &subscribe); err != nil {
						ctr.Logger.Error(ctx, "failed to handle SubscribeResponseMessage",
							logger.Value("error", err))
						done <- TerminateTypeSubscribeMsgHandlerError
						return
					}
				}
			case ResponseTypeUnsubscribeResponse:
				var unsubscribe UnsubscribeResponseMessage
				err := utils.UnmarshalJSON(content, &unsubscribe)
				if err != nil {
					ctr.Logger.Error(ctx, "failed to read UnsubscribeResponseMessage",
						logger.Value("error", err))
					if conf.DisconnectOnParseMsgError {
						done <- TerminateTypeParseMsgError
						return
					}
				}
				if v, ok := ws.unsubscribeMsgMem.Memory[unsubscribe.MessageID]; ok {
					if unsubscribe.Success {
						for _, id := range v.SubscriptionIDs {
							delete(ws.SubscribeMemory.Memory, id)

							if consumerID, ok := ws.SubscribeConsumerMemorySubIndex.Memory[id]; ok {
								delete(ws.SubscribeConsumerMemory.Memory, consumerID)
								delete(ws.SubscribeConsumerMemorySubIndex.Memory, id)
							}
						}
						delete(ws.unsubscribeMsgMem.Memory, unsubscribe.MessageID)
					}
				}
				if ws.UnsubscribeMsgHandler != nil {
					if err := ws.UnsubscribeMsgHandler(ws, &unsubscribe); err != nil {
						ctr.Logger.Error(ctx, "failed to handle UnsubscribeResponseMessage",
							logger.Value("error", err))
						done <- TerminateTypeUnsubscribeMsgHandlerError
						return
					}
				}
			case ResponseTypeUnsupportedResponse:
				var unsupported UnsupportedResponseMessage
				err := utils.UnmarshalJSON(content, &unsupported)
				if err != nil {
					ctr.Logger.Error(ctx, "failed to read UnsupportedResponseMessage",
						logger.Value("error", err))
					if conf.DisconnectOnParseMsgError {
						done <- TerminateTypeParseMsgError
						return
					}
				}
				if ws.UnsupportedMsgHandler != nil {
					if err := ws.UnsupportedMsgHandler(ws, &unsupported); err != nil {
						ctr.Logger.Error(ctx, "failed to handle UnsupportedResponseMessage",
							logger.Value("error", err))
						done <- TerminateTypeUnsupportedMsgHandlerError
						return
					}
				}
			case ResponseTypeEventResponse:
				var event EventResponseMessage
				err := utils.UnmarshalJSON(content, &event)
				if err != nil {
					ctr.Logger.Error(ctx, "failed to read EventResponseMessage",
						logger.Value("error", err))
					if conf.DisconnectOnParseMsgError {
						done <- TerminateTypeParseMsgError
						return
					}
				}
				// fmt.Println("EventResponseMessage", event)
				if ws.EventMsgHandler != nil {
					if ws.EventMsgHandler(ws, &event, content); err != nil {
						ctr.Logger.Error(ctx, "failed to handle EventResponseMessage",
							logger.Value("error", err))
						done <- TerminateTypeEventMsgHandlerError
						return
					}
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

func (ws *WebSocket) SendSubscribeMessage(consumerId string, aggregateType AggregateType, aggregateIDs []string, eventTypes []EventType) error {
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
	ws.subscribeMsgMem.Memory[msgID] = SubscribeMessageMemoryData{ConsumerID: consumerId, AggregateType: aggregateType, AggregateIDs: aggregateIDs, EventTypes: eventTypes}
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

func (ws *WebSocket) UnsubscribeByConsumerID(consumerID string) error {
	var subscriptionID string
	if id, ok := ws.SubscribeConsumerMemory.Memory[consumerID]; ok {
		subscriptionID = id
	} else {
		return nil
	}

	return ws.SendUnsubscribeMessage([]string{subscriptionID})
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
