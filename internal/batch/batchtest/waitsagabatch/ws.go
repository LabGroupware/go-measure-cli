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

func NewSocket() *Socket {
	sock := ws.NewWebSocket()

	sock.CannotParseMsgHandler = ws.DefaultCannotParseResponseMessageHandleFunc
	sock.SubscribeMsgHandler = ws.DefaultSubscribeResponseMessageHandleFunc
	sock.UnsubscribeMsgHandler = ws.DefaultUnsubscribeResponseMessageHandleFunc
	sock.UnsupportedMsgHandler = ws.DefaultUnsupportedResponseMessageHandleFunc
	sock.EventMsgHandler = ws.DefaultEventResponseMessageHandleFunc

	return &Socket{
		ws: sock,
	}
}

func (s *Socket) Connect(ctx context.Context, ctr *app.Container) (<-chan struct{}, error) {
	done, err := s.ws.Connect(ctx, ctr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %v", err)
	}

	return done, nil
}
