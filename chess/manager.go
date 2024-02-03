package chess

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	// Map from gameID to *Game
	games map[string]*Game
	// Map from userID to gameID
	players map[string]string
	// Lock when reading or writing
	sync.RWMutex
}

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  256,
		WriteBufferSize: 8192,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func NewManager() *Manager {
	return &Manager{
		games:   make(map[string]*Game),
		players: make(map[string]string),
	}
}

func (m *Manager) MakeGame(playerIDs []string, playerConns []*websocket.Conn) {
	m.Lock()
	defer m.Unlock()

	// Make the game here
}

func (m *Manager) RemoveGame(game *Game) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.games[game.GameID]; ok {
		// close connection
		// client.connection.Close()
		// remove
		delete(m.games, game.GameID)
	}
}
