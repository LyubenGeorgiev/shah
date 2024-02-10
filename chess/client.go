package chess

import (
	"bytes"
	"context"
	"fmt"
	"time"

	boardview "github.com/LyubenGeorgiev/shah/view/board"
	"github.com/LyubenGeorgiev/shah/view/board/models"
	"github.com/gorilla/websocket"
)

const (
	// pongWait is how long we will await a pong response from client
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type Output struct {
	Type    string
	Payload interface{}
}

type Client struct {
	game    *Game
	conn    *websocket.Conn
	outputs chan Output
	Side    Side
	ctx     context.Context
}

func NewClient(game *Game, conn *websocket.Conn, side Side, ctx context.Context) *Client {
	return &Client{
		game:    game,
		conn:    conn,
		outputs: make(chan Output),
		Side:    side,
		ctx:     ctx,
	}
}

func (c *Client) ListenInput() {
	defer func() {
		// Graceful Close the Connection once this
		c.game.removeClient(c)
	}()

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		fmt.Println(err)
		return
	}
	// Configure how to handle Pong responses
	c.conn.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error reading message: %v\n", err)
			}
			break
		}

		inputEvent, err := UnmarshalAction(payload)
		if err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		inputEvent.Client = c

		c.game.InputEvent(inputEvent)
	}
}

// pongHandler is used to handle PongMessages for the Client
func (c *Client) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	return c.conn.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) ListenOutput() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		// Graceful close if this triggers a closing
		c.game.removeClient(c)
	}()

	for {
		select {
		case output, ok := <-c.outputs:
			if !ok {
				fmt.Println("Outputs channel is closed.")
				return
			}

			var buf bytes.Buffer
			if output.Type == "board" {
				bs, ok := output.Payload.(*models.BoardState)
				if !ok {
					continue
				}

				if err := boardview.Board(bs).Render(c.ctx, &buf); err != nil {
					fmt.Printf("Error rendering board: %v\n", err)
					continue
				}
			} else if output.Type == "chat" {
				messages, ok := output.Payload.([]models.Message)
				if !ok {
					continue
				}

				if err := boardview.Chat(messages, c.game.getID(c)).Render(c.ctx, &buf); err != nil {
					fmt.Printf("Error rendering chat: %v\n", err)
					continue
				}
			} else if output.Type == "message" {
				message, ok := output.Payload.(models.Message)
				if !ok {
					continue
				}

				if err := boardview.Message(message, c.game.getID(c)).Render(c.ctx, &buf); err != nil {
					fmt.Printf("Error rendering chat: %v\n", err)
					continue
				}
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
				fmt.Printf("Error writing to websocket: %v\n", err)
				continue
			}
		case <-ticker.C:
			// Send the Ping
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				fmt.Println("Ping error: ", err)
				return
			}
		}
	}
}
