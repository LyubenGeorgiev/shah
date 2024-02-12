package application

import (
	"net/http"
	"strconv"

	"github.com/LyubenGeorgiev/shah/chess"
	"github.com/LyubenGeorgiev/shah/view/replay"
	"github.com/gorilla/mux"
)

func (app *App) HandleReplays(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)["gameID"]
	game, err := app.Storage.GetGame(gameID)
	if err != nil {
		http.Error(w, "Game not found!", http.StatusNotFound)
		return
	}

	replay.ReplayLayout(&replay.BoardState{Pieces: chess.NewBoadFromFen(chess.STARTPOS_FEN).GetPieces()}, gameID, game.WhiteID, game.BlackID, 0, len(game.Moves)).Render(r.Context(), w)
}

func (app *App) HandleReplaysMove(w http.ResponseWriter, r *http.Request) {
	gameID, movesPlayedStr := mux.Vars(r)["gameID"], mux.Vars(r)["move"]
	movesPlayed, err := strconv.ParseInt(movesPlayedStr, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	game, err := app.Storage.GetGame(gameID)
	if err != nil {
		http.Error(w, "Game not found!", http.StatusNotFound)
		return
	}

	board := chess.Startpos()
	for i := 0; i < int(movesPlayed); i++ {
		board.MakeMove(chess.Move(game.Moves[i]), false)
	}

	replay.Replay(&replay.BoardState{Pieces: board.GetPieces()}, gameID, game.WhiteID, game.BlackID, int(movesPlayed), len(game.Moves)).Render(r.Context(), w)
}
