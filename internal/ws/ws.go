package ws

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

// WebSocket is a struct that represents a WebSocket connection
type WebSocket struct {
	conn *websocket.Conn
}

// NewWebSocket creates a new WebSocket connection
func NewWebSocket() *WebSocket {
	return &WebSocket{}
}

// Connect connects to a remote server using the WebSocket protocol
func (ws *WebSocket) Connect(addr string) (<-chan struct{}, error) {
	headers := http.Header{}
	conn, _, err := websocket.DefaultDialer.Dial(addr, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket server: %w", err)
	}
	ws.conn = conn

	// Clean up the connection when the function returns
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			var msg Message
			err := ws.conn.ReadJSON(&msg)
			if err != nil {
				fmt.Printf("failed to read message: %v\n", err)
				return
			}
			fmt.Printf("received message: %v\n", msg)
		}
	}()

	return done, nil
}

// Send sends a message to the server
func (ws *WebSocket) Send(message string) error {
	err := ws.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
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
