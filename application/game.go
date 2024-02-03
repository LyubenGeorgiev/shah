package application

import (
	"context"
	"fmt"
	"log"
	"net/http"

	// "github.com/LyubenGeorgiev/shah/view/components"
	// "github.com/LyubenGeorgiev/shah/view/components/models"

	"github.com/LyubenGeorgiev/shah/chess"
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

	// sig   = make(chan struct{})
	ids   = []string{}
	conns = []*websocket.Conn{}
)

// TODO make actual handler :D
func (a *App) Game(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil || userID == "" {
		log.Println("Unknown user!", userID, "Error", err.Error())
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}
	fmt.Println("Request", userID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	ids = append(ids, userID)
	conns = append(conns, conn)

	if len(conns) == 2 {
		go chess.NewGame(ids[0], ids[1], conns[0], conns[1], context.Background(), context.Background()).StartGame()
	}
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
		layout.Play().Render(r.Context(), w)

		return
	} else if err != nil {
		log.Println("Failed to get game for userID", userID, "Error:", err.Error())
		http.Error(w, "Failed getting game", http.StatusInternalServerError)
		return
	}

	// Here we know user is ingame TODO replace this
	fmt.Println("TODO use gameid", gameID)
}
