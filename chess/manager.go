package chess

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/LyubenGeorgiev/shah/db"
	"github.com/LyubenGeorgiev/shah/db/models"
	"github.com/LyubenGeorgiev/shah/util"
	boardview "github.com/LyubenGeorgiev/shah/view/board"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Manager struct {
	// Map from gameID to *Game
	games map[string]*Game
	// Map from userID to gameID
	players map[string]string
	// Lock when reading or writing
	sync.RWMutex

	// Queue for matchmaking
	queue chan string

	Storage db.Storage
}

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  256,
		WriteBufferSize: 8192,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewManager(storage db.Storage) *Manager {
	return &Manager{
		games:   make(map[string]*Game),
		players: make(map[string]string),
		queue:   make(chan string),
		Storage: storage,
	}
}

func (m *Manager) HandleGame(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["id"]

	userID, err := util.GetUserID(r)
	if err != nil || userID == "" {
		log.Println("Unknown user!", userID, "Error", err.Error())
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	// User is trying to access different game than his own
	if gameID != m.GetUserGameID(userID) {
		log.Println("Unknown user!", userID, "Error", err.Error())
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	// Here we know user is ingame so we need to connect him to the game
	if game := m.GetGame(gameID); game != nil {
		if err := game.Connect(w, r, userID); err == nil {
			// All good return
			return
		} else {
			log.Println(err)
		}
	}

	http.Error(w, "Error connecting user to a game!", http.StatusInternalServerError)
}

func (m *Manager) HandlePlay(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserID(r)
	if err != nil || userID == "" {
		log.Println("Unknown user!", userID, "Error", err.Error())
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	if gameID := m.GetUserGameID(userID); gameID != "" {
		boardview.BoardWebsocket(gameID).Render(r.Context(), w)
		return
	}

	// Add user to queue
	select {
	case m.queue <- userID:
	case opponentUserID := <-m.queue:
		m.MakeGame(opponentUserID, userID)
	}

	if gameID := m.GetUserGameID(userID); gameID != "" {
		boardview.BoardWebsocket(gameID).Render(r.Context(), w)
		return
	}
}

func (m *Manager) GetUserGameID(userID string) string {
	m.Lock()
	defer m.Unlock()

	if gameID, ok := m.players[userID]; ok {
		return gameID
	}

	return ""
}

func (m *Manager) GetGame(gameID string) *Game {
	m.Lock()
	defer m.Unlock()

	if game, ok := m.games[gameID]; ok {
		return game
	}

	return nil
}

func (m *Manager) MakeGame(whiteID, blackID string) {
	m.Lock()
	defer m.Unlock()

	game := NewGame(m, whiteID, blackID, 10*time.Minute)
	m.games[game.GameID] = game
	m.players[whiteID] = game.GameID
	m.players[blackID] = game.GameID

	go game.Start()
}

func (m *Manager) RemoveGame(game *Game, winnerID string) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.games[game.GameID]; ok {
		if winnerID != "" {
			m.Storage.CreateGame(&models.Game{ID: game.GameID, WhiteID: game.whiteID, BlackID: game.blackID, WinnerID: sql.NullString{String: winnerID, Valid: true}, Moves: game.Moves})
		} else {
			m.Storage.CreateGame(&models.Game{ID: game.GameID, WhiteID: game.whiteID, BlackID: game.blackID, Moves: game.Moves})
		}

		delete(m.players, game.whiteID)
		delete(m.players, game.blackID)
		delete(m.games, game.GameID)
	}
}
