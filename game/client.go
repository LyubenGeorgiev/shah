package game

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/LyubenGeorgiev/shah/view/components"
	"github.com/LyubenGeorgiev/shah/view/components/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	userID        string
	conn          *websocket.Conn
	ctx           context.Context
	inputs        chan<- *inputEvent
	outputs       <-chan *models.BoardState
	side          models.Side
	remainingTime time.Time
}

func (c *Client) ListenInput() {
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
		fmt.Printf("%+v\n", inputEvent)

		c.inputs <- inputEvent
	}
}

func (c *Client) ListenOutput() {
	for {
		select {
		case <-c.ctx.Done():
			// TODO handle context done
		case bs, ok := <-c.outputs:
			if !ok {
				fmt.Println("Outputs channel is closed.")
				return
			}

			var buf bytes.Buffer
			if err := components.Board(bs).Render(c.ctx, &buf); err != nil {
				fmt.Printf("Error rendering board: %v\n", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
				fmt.Printf("Error writing to websocket: %v\n", err)
				continue
			}
		}
	}
}
