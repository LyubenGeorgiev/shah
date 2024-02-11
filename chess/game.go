package chess

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/LyubenGeorgiev/shah/view/board/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Game struct {
	Manager            *Manager
	GameID             string
	whiteID            string
	blackID            string
	whiteRemainingTime time.Duration
	blackRemainingTime time.Duration
	board              Board
	lastMove           *Move
	inputs             chan *inputEvent
	white              *Client
	black              *Client
	Moves              pq.Int32Array
	Messages           []models.Message
}

func NewGame(manager *Manager, whiteID, blackID string, timeGiven time.Duration) *Game {
	id := uuid.Must(uuid.NewRandom()).String()
	g := &Game{
		Manager:            manager,
		GameID:             id,
		whiteID:            whiteID,
		blackID:            blackID,
		whiteRemainingTime: timeGiven,
		blackRemainingTime: timeGiven,
		board:              *NewBoadFromFen([]byte("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 ")),
		lastMove:           nil,
		inputs:             make(chan *inputEvent),
		white:              nil,
		black:              nil,
		Moves:              pq.Int32Array{},
		Messages:           []models.Message{},
	}

	return g
}

func (g *Game) InputEvent(ie *inputEvent) {
	g.inputs <- ie
}

func (g *Game) removeClient(c *Client) {
	c.conn.Close()
	if g.white == c {
		g.white = nil
	} else if g.black == c {
		g.black = nil
	}
}

func (g *Game) Connect(w http.ResponseWriter, r *http.Request, userID string) error {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	if g.whiteID == userID {
		if g.white != nil {
			g.removeClient(g.white)
		}

		g.white = NewClient(g, conn, white, r.Context())

		go g.white.ListenInput()
		go g.white.ListenOutput()

		g.white.outputs <- Output{Type: "board", Payload: &models.BoardState{
			Highlighted: nil,
			Pieces:      g.getPieces(),
			Moves:       nil,
			Captures:    nil,
			View:        models.Side(white),
		}}
	} else if g.blackID == userID {
		if g.black != nil {
			g.removeClient(g.black)
		}

		g.black = NewClient(g, conn, black, r.Context())

		go g.black.ListenInput()
		go g.black.ListenOutput()

		g.black.outputs <- Output{Type: "board", Payload: &models.BoardState{
			Highlighted: nil,
			Pieces:      g.getPieces(),
			Moves:       nil,
			Captures:    nil,
			View:        models.Side(black),
		}}
	}

	return nil
}

func (g *Game) Start() {
	moves := make(chan Move)
	for g.board.Gameover() {
		timer := g.getTimer()

		ctx, cancel := context.WithCancel(context.Background())

		go g.handleInputEvents(ctx, moves)

		select {
		case move := <-moves:
			g.Moves = append(g.Moves, int32(move))
		case <-timer.C:
			cancel()
			close(moves)
			break
		}
		cancel()
		timer.Stop()
	}

	king := k
	if g.board.Side == white {
		king = K
	}

	// If king is in check current side lost otherwise it is a draw
	if g.board.isSquareAttacked(square(g.board.Bitboards[king].GetLs1bIndex()), g.board.Side.opposite()) {
		if g.board.Side == white {
			g.Manager.RemoveGame(g, g.blackID)
		} else {
			g.Manager.RemoveGame(g, g.whiteID)
		}
	} else {
		g.Manager.RemoveGame(g, "")
	}
}

func (g *Game) getTimer() *time.Timer {
	if g.board.Side == white {
		return time.NewTimer(g.whiteRemainingTime)
	}

	return time.NewTimer(g.blackRemainingTime)
}

func (g *Game) getID(c *Client) string {
	if g.white == c {
		return g.whiteID
	}

	return g.blackID
}

func (g *Game) handleInputEvents(ctx context.Context, moves chan<- Move) {
	var selected square
	var legalMoves *Moves

	for {
		select {
		case event, ok := <-g.inputs:
			if !ok {
				fmt.Println("Input event channel was closed unexpectedly!")
			}

			if event.Action == "chat" {
				UpdateClient(event.Client, Output{Type: event.Action, Payload: g.Messages})
				continue
			} else if event.Action == "message" {
				msg := models.Message{Text: event.Message, UserID: g.getID(event.Client)}
				UpdateClient(g.white, Output{Type: event.Action, Payload: msg})
				UpdateClient(g.black, Output{Type: event.Action, Payload: msg})
				g.Messages = append(g.Messages, msg)
				continue
			}

			// No events for opposite side player
			if event.Client.Side != g.board.Side {
				continue
			}

			if event.Action == "unselect" {
				selected = no_sq

				// Highlight selected and last move
				highlighted := []int{}
				if g.lastMove != nil {
					highlighted = append(highlighted, int(g.lastMove.getSource()))
					highlighted = append(highlighted, int(g.lastMove.getTarget()))
				}

				UpdateClient(event.Client, Output{Type: "board", Payload: &models.BoardState{
					Highlighted: highlighted,
					Pieces:      g.getPieces(),
					Moves:       nil,
					Captures:    nil,
					View:        models.Side(event.Client.Side),
				}})
			} else if event.Action == "select" {
				selected = stringToSquare[event.Square]
				legalMoves = g.board.GetLegalMoves(false).FilterSelected(selected)

				// Highlight selected and last move
				highlighted := []int{int(stringToSquare[event.Square])}
				if g.lastMove != nil {
					highlighted = append(highlighted, int(g.lastMove.getSource()))
					highlighted = append(highlighted, int(g.lastMove.getTarget()))
				}

				// Highlight possible moves
				quietMoves := []int{}
				captureMoves := []int{}
				for i := 0; i < legalMoves.count; i++ {
					if legalMoves.moves[i].isCapture() {
						captureMoves = append(captureMoves, int(legalMoves.moves[i].getTarget()))
					} else {
						quietMoves = append(quietMoves, int(legalMoves.moves[i].getTarget()))
					}
				}

				UpdateClient(event.Client, Output{Type: "board", Payload: &models.BoardState{
					Highlighted: highlighted,
					Pieces:      g.getPieces(),
					Moves:       quietMoves,
					Captures:    captureMoves,
					View:        models.Side(event.Client.Side),
				}})
			} else if event.Action == "move" {
				target := stringToSquare[event.Square]
				legalMoves = g.board.GetLegalMoves(false).FilterSelected(selected)

				for i := 0; i < legalMoves.count; i++ {
					if legalMoves.moves[i].getTarget() == target {
						g.board.makeMove(legalMoves.moves[i], false)
						g.lastMove = &legalMoves.moves[i]
						go func(m Move, moves chan<- Move) {
							moves <- m
						}(legalMoves.moves[i], moves)
						break
					}
				}

				// Highlight selected and last move
				highlighted := []int{}
				if g.lastMove != nil {
					highlighted = append(highlighted, int(g.lastMove.getSource()))
					highlighted = append(highlighted, int(g.lastMove.getTarget()))
				}

				UpdateClient(g.white, Output{Type: "board", Payload: &models.BoardState{
					Highlighted: highlighted,
					Pieces:      g.getPieces(),
					Moves:       nil,
					Captures:    nil,
					View:        models.Side(white),
				}})

				UpdateClient(g.black, Output{Type: "board", Payload: &models.BoardState{
					Highlighted: highlighted,
					Pieces:      g.getPieces(),
					Moves:       nil,
					Captures:    nil,
					View:        models.Side(black),
				}})

				return
			} else {
				fmt.Println("Unknown action:", event.Action)
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

func (g *Game) getPieces() map[int]string {
	pieces := map[int]string{}

	for bb := P; bb <= k; bb++ {
		bitboard := g.board.Bitboards[bb]

		for bitboard > 0 {
			sourceSquare := square(bitboard.GetLs1bIndex())

			pieces[int(sourceSquare)] = string(pieceToChar[bb])

			bitboard.popBit(sourceSquare)
		}
	}

	return pieces
}

func UpdateClient(c *Client, update Output) {
	if c != nil {
		c.outputs <- update
	}
}
