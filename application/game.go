package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/LyubenGeorgiev/shah/view/components"
	"github.com/LyubenGeorgiev/shah/view/components/models"
	"github.com/LyubenGeorgiev/shah/view/layout"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Game struct {
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  256,
		WriteBufferSize: 8192,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	first = true
)

type Action struct {
	Square string `json:"HX-Trigger"`
	Action string `json:"HX-Trigger-Name"`
	URL    string `json:"HX-Current-URL"`
}

func UnmarshalAction(bytes []byte) (*Action, error) {
	var tmp struct {
		Headers Action `json:"HEADERS"`
	}

	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return nil, err
	}

	return &tmp.Headers, nil
}

// TODO make actual handler :D
func (a *App) Game(w http.ResponseWriter, r *http.Request) {
	if !first {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			action, err := UnmarshalAction(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%+v\n", action)
			fmt.Println("Message:", string(msg), "Type:", msgType, "Size:", len(msg))

			var buf bytes.Buffer

			err = layout.BoardTest(&models.BoardState{
				Highlighted: []int{},
				Pieces: map[int]string{
					0: "r", 1: "n", 2: "b", 3: "q", 4: "k", 5: "b", 6: "n", 7: "r",
					8: "p", 9: "p", 10: "p", 11: "p", 12: "p", 13: "p", 14: "p", 15: "p",
					48: "P", 49: "P", 50: "P", 51: "P", 52: "P", 53: "P", 54: "P", 55: "P",
					56: "R", 57: "N", 58: "B", 59: "Q", 60: "K", 61: "B", 62: "N", 63: "R",
				},
				Moves:    []int{},
				Captures: []int{},
				View:     models.White,
			}).Render(r.Context(), &buf)
			if err != nil {
				fmt.Printf("Error rendering game board: %v", err)
			}

			fmt.Println(len(buf.Bytes()))
			// conn.WriteMessage(websocket.TextMessage, buf.Bytes())
		}
	}

	err := layout.BoardTest(&models.BoardState{
		Highlighted: []int{},
		Pieces: map[int]string{
			0: "r", 1: "n", 2: "b", 3: "q", 4: "k", 5: "b", 6: "n", 7: "r",
			8: "p", 9: "p", 10: "p", 11: "p", 12: "p", 13: "p", 14: "p", 15: "p",
			48: "P", 49: "P", 50: "P", 51: "P", 52: "P", 53: "P", 54: "P", 55: "P",
			56: "R", 57: "N", 58: "B", 59: "Q", 60: "K", 61: "B", 62: "N", 63: "R",
		},
		Moves:    []int{},
		Captures: []int{},
		Unselect: []int{8,
			16, 17, 18, 19, 20, 21, 22, 23,
			24, 25, 26, 27, 28, 29, 30, 31,
			32, 33, 34, 35, 36, 37, 38, 39,
			40, 41, 42, 43, 44, 45, 46, 47,
		},
		View: models.White,
	}).Render(r.Context(), w)
	if err != nil {
		fmt.Printf("Error rendering game board: %v", err)
	}

	first = false
}

func (a *App) Play(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil || userID == "" {
		log.Println("Unknown user!", userID, "Error", err.Error())
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	gameID, err := a.Cache.UserIsIngame(r.Context(), userID)
	if err == redis.Nil {
		// Add this user to the queue for matchmaking
		if gameID != "" {

		}
	} else if err != nil {
		log.Println("Failed to get game for userID", userID, "Error:", err.Error())
		http.Error(w, "Failed getting game", http.StatusInternalServerError)
		return
	}

	// Here we know user is ingame TODO replace this

	layout.Play().Render(r.Context(), w)
}
