package application

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LyubenGeorgiev/shah/chess"
	"github.com/LyubenGeorgiev/shah/util"
	"github.com/LyubenGeorgiev/shah/view/computer"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *App) ComputerGame(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserID(r)
	if err != nil || userID == "" {
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	gameID := mux.Vars(r)["gameID"]
	if cacheGameID, err := a.Cache.GetUserInComputerGame(r.Context(), userID); err != nil || cacheGameID != gameID {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	strMoves, _ := a.Cache.GetComputerGamestate(r.Context(), gameID)
	board := a.getBoard(r.Context(), gameID)

	action, source := r.Header.Get("HX-Trigger-Name"), r.Header.Get("HX-Trigger")

	if action == "unselect" {
		// Highlight selected and last move
		highlighted := []int{}
		if len(strMoves) > 0 {
			highlighted = append(highlighted, parseMove(strMoves[len(strMoves)-1]).GetSource())
			highlighted = append(highlighted, parseMove(strMoves[len(strMoves)-1]).GetTarget())
		}

		computer.ComputerBoard(&computer.BoardState{
			GameID:      gameID,
			Highlighted: highlighted,
			Pieces:      board.GetPieces(),
			Moves:       nil,
			Captures:    nil,
			View:        computer.White,
		}).Render(r.Context(), w)
	} else if action == "select" {
		selected := chess.StringToSquare[source]
		legalMoves := board.GetLegalMoves(false).FilterSelected(selected)

		// Highlight selected and last move
		highlighted := []int{chess.StringToSquare[source]}
		if len(strMoves) > 0 {
			highlighted = append(highlighted, parseMove(strMoves[len(strMoves)-1]).GetSource())
			highlighted = append(highlighted, parseMove(strMoves[len(strMoves)-1]).GetTarget())
		}

		// Highlight possible moves
		quietMoves := map[int]string{}
		captureMoves := map[int]string{}
		for i := 0; i < legalMoves.Count(); i++ {
			if legalMoves.At(i).IsCapture() {
				captureMoves[legalMoves.At(i).GetTarget()] = fmt.Sprintf("%d", legalMoves.At(i))
			} else {
				quietMoves[legalMoves.At(i).GetTarget()] = fmt.Sprintf("%d", legalMoves.At(i))
			}
		}

		computer.ComputerBoard(&computer.BoardState{
			GameID:      gameID,
			Highlighted: highlighted,
			Pieces:      board.GetPieces(),
			Moves:       quietMoves,
			Captures:    captureMoves,
			View:        computer.White,
		}).Render(r.Context(), w)
	} else { // We are making a move
		move := parseMove(action)
		if !board.MakeMove(move, false) {
			http.Error(w, "Illegal move", http.StatusInternalServerError)
			return
		}

		a.Cache.PushComputerGamestateMove(r.Context(), gameID, action)

		// Handle gameover
		if board.Gameover() {
			a.Cache.DelGamestate(r.Context(), gameID)
			a.Cache.DelUserInComputerGame(r.Context(), userID)
			fmt.Println("User won vs bot")
			return
		}

		eng := chess.NewEngine()
		eng.Board = board
		botMove := eng.Search()

		board.MakeMove(botMove, false)
		a.Cache.PushComputerGamestateMove(r.Context(), gameID, fmt.Sprintf("%d", botMove))

		// Handle gameover
		if board.Gameover() {
			a.Cache.DelGamestate(r.Context(), gameID)
			a.Cache.DelUserInComputerGame(r.Context(), userID)
			fmt.Println("Bot won vs user")
			return
		}

		computer.ComputerBoard(&computer.BoardState{
			GameID:      gameID,
			Highlighted: nil,
			Pieces:      board.GetPieces(),
			Moves:       nil,
			Captures:    nil,
			View:        computer.White,
		}).Render(r.Context(), w)
	}
}

func (a *App) Computer(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserID(r)
	if err != nil || userID == "" {
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	gameID, err := a.Cache.GetUserInComputerGame(r.Context(), userID)
	if err != nil {
		gameID = uuid.Must(uuid.NewRandom()).String()
		a.Cache.SetUserInComputerGame(r.Context(), userID, gameID)
		a.Cache.SetComputerGamestate(r.Context(), gameID)
	}

	board := a.getBoard(r.Context(), gameID)

	computer.Layout(userID, &computer.BoardState{
		GameID:      gameID,
		Highlighted: nil,
		Pieces:      board.GetPieces(),
		Moves:       nil,
		Captures:    nil,
		View:        computer.White,
	}).Render(r.Context(), w)
}

func parseMove(str string) chess.Move {
	move, _ := strconv.ParseInt(str, 10, 32)
	return chess.Move(move)
}

func (a *App) getBoard(ctx context.Context, gameID string) chess.Board {
	moves, _ := a.Cache.GetComputerGamestate(ctx, gameID)

	board := chess.Startpos()
	for _, move := range moves {
		board.MakeMove(parseMove(move), false)
	}

	return board
}
