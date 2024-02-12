package chess

import (
	"database/sql"
	"log"
	"math"
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
		game := m.GetGame(gameID)
		if game.whiteID == userID {
			boardview.BoardWebsocket(gameID, game.whiteID, game.blackID).Render(r.Context(), w)
		} else {
			boardview.BoardWebsocket(gameID, game.blackID, game.whiteID).Render(r.Context(), w)
		}
		return
	}

	// Add user to queue
	select {
	case m.queue <- userID:
	case opponentUserID := <-m.queue:
		m.MakeGame(opponentUserID, userID)
	}

	if gameID := m.GetUserGameID(userID); gameID != "" {
		game := m.GetGame(gameID)
		if game.whiteID == userID {
			boardview.BoardWebsocket(gameID, game.whiteID, game.blackID).Render(r.Context(), w)
		} else {
			boardview.BoardWebsocket(gameID, game.blackID, game.whiteID).Render(r.Context(), w)
		}
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
			if winnerID == game.whiteID {
				m.UpdateRating(game.whiteID, game.blackID, false)
			} else {
				m.UpdateRating(game.blackID, game.whiteID, false)
			}
		} else {
			m.Storage.CreateGame(&models.Game{ID: game.GameID, WhiteID: game.whiteID, BlackID: game.blackID, Moves: game.Moves})
			m.UpdateRating(game.whiteID, game.blackID, true)
		}

		delete(m.players, game.whiteID)
		delete(m.players, game.blackID)
		delete(m.games, game.GameID)
	}
}

func (m *Manager) UpdateRating(winnerID, loserID string, isDraw bool) {
	winner, _ := m.Storage.FindByUserID(winnerID)
	loser, _ := m.Storage.FindByUserID(loserID)

	KFactor := 32.0

	// Calculate expected scores
	expectedScore1 := 1.0 / (1.0 + math.Pow(10, (loser.Rating-winner.Rating)/400.0))
	expectedScore2 := 1.0 / (1.0 + math.Pow(10, (winner.Rating-loser.Rating)/400.0))

	// If it's a draw, adjust ratings based on the expected scores
	if isDraw {
		// Adjust rating for player 1
		winner.Rating += KFactor * (0.5 - expectedScore1)
		// Adjust rating for player 2
		loser.Rating += KFactor * (0.5 - expectedScore2)
	} else {
		// Calculate the change in ratings for a decisive result
		delta1 := KFactor * (1 - expectedScore1)
		delta2 := KFactor * (0 - expectedScore2)

		// Update ratings for players based on the result
		winner.Rating += delta1
		loser.Rating += delta2

		winner.GamesWon++
	}

	// Update games played for both players
	winner.GamesPlayed++
	loser.GamesPlayed++

	m.Storage.SaveUser(winner)
	m.Storage.SaveUser(loser)
}
