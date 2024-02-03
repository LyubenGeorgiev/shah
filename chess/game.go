package chess

import (
	"context"
	"fmt"
	"time"

	"github.com/LyubenGeorgiev/shah/view/components/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Game struct {
	GameID   string
	board    Board
	lastMove *Move
	inputs   chan *inputEvent
	white    *Client
	black    *Client
}

func NewGame(whiteID, blackID string, whiteConn, blackConn *websocket.Conn, whiteCtx, blackCtx context.Context) *Game {
	id := uuid.Must(uuid.NewRandom()).String()
	g := &Game{
		GameID:   id,
		board:    *NewBoadFromFen([]byte("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 ")),
		lastMove: nil,
		inputs:   make(chan *inputEvent),
	}

	g.white = NewClient(whiteID, g, whiteConn, make(chan *models.BoardState), white, 10*time.Minute, whiteCtx)
	g.black = NewClient(blackID, g, blackConn, make(chan *models.BoardState), black, 10*time.Minute, blackCtx)

	return g
}

func (g *Game) InputEvent(ie *inputEvent) {
	g.inputs <- ie
}

func (g *Game) removeClient(c *Client) {
	c.conn.Close()
	if g.white == c {
		g.white = nil
	} else {
		g.black = nil
	}
}

func (g *Game) StartGame() {
	go g.white.ListenInput()
	go g.white.ListenOutput()
	go g.black.ListenInput()
	go g.black.ListenOutput()

	g.white.outputs <- &models.BoardState{
		Highlighted: nil,
		Pieces:      g.getPieces(),
		Moves:       nil,
		Captures:    nil,
		View:        models.Side(white),
	}

	g.black.outputs <- &models.BoardState{
		Highlighted: nil,
		Pieces:      g.getPieces(),
		Moves:       nil,
		Captures:    nil,
		View:        models.Side(black),
	}

	moves := make(chan Move)
	for g.board.Gameover() {
		timer := g.getTimer()

		ctx, cancel := context.WithCancel(context.Background())

		go g.handleInputEvents(ctx, moves)

		select {
		case move := <-moves:
			fmt.Println(move)
		case <-timer.C:
			cancel()
			close(moves)
			break
		}
		cancel()
		timer.Stop()
	}
	// Game is over handle the winner
}

func (g *Game) getTimer() *time.Timer {
	if g.board.Side == white {
		return time.NewTimer(g.white.RemainingTime)
	}

	return time.NewTimer(g.black.RemainingTime)
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

				event.Client.outputs <- &models.BoardState{
					Highlighted: highlighted,
					Pieces:      g.getPieces(),
					Moves:       nil,
					Captures:    nil,
					View:        models.Side(event.Client.Side),
				}
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

				event.Client.outputs <- &models.BoardState{
					Highlighted: highlighted,
					Pieces:      g.getPieces(),
					Moves:       quietMoves,
					Captures:    captureMoves,
					View:        models.Side(event.Client.Side),
				}
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

				g.white.outputs <- &models.BoardState{
					Highlighted: highlighted,
					Pieces:      g.getPieces(),
					Moves:       nil,
					Captures:    nil,
					View:        models.Side(white),
				}

				g.black.outputs <- &models.BoardState{
					Highlighted: highlighted,
					Pieces:      g.getPieces(),
					Moves:       nil,
					Captures:    nil,
					View:        models.Side(black),
				}

				return
			} else {
				fmt.Println("Unknown action!")
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
