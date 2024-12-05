package waitsagabatch

import (
	"context"
	"fmt"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/ws"
)

type Socket struct {
	ws *ws.WebSocket
}

func NewSocket(ctx context.Context, ctr *app.Container) (*Socket, error) {
	sock := ws.NewWebSocket()

	sock.CannotParseMsgHandler = func(ws *ws.WebSocket, msg *ws.CannotParseResponseMessage) {
		fmt.Printf("cannot parse response: %s\n", msg.MessageID)
	}
	sock.SubscribeMsgHandler = func(ws *ws.WebSocket, msg *ws.SubscribeResponseMessage) {
		fmt.Printf("subscribed: %s\n", msg.SubscriptionID)
	}
	sock.UnsubscribeMsgHandler = func(ws *ws.WebSocket, msg *ws.UnsubscribeResponseMessage) {
		fmt.Printf("unsubscribed: %s\n", msg.MessageID)
	}
	sock.UnsupportedMsgHandler = func(ws *ws.WebSocket, msg *ws.UnsupportedResponseMessage) {
		fmt.Printf("unsupported response: %s\n", msg.MessageID)
	}
	sock.EventMsgHandler = func(ws *ws.WebSocket, msg *ws.EventResponseMessage, raw []byte) {
		fmt.Printf("event response: %s\n", msg.EventType)
	}

	return &Socket{
		ws: sock,
	}, nil
}

func (s *Socket) Connect(ctx context.Context, ctr *app.Container) (<-chan struct{}, error) {
	done, err := s.ws.Connect(ctx, ctr, ws.ConnectConfig{
		DisconnectOnReadMsgError:       true,
		DisconnectOnUnmarshalJSONError: true,
		DisconnectOnParseMsgError:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %v", err)
	}

	return done, nil
}
