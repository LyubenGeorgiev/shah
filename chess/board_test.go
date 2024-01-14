package chess

import (
	"fmt"
	"testing"
)

func (m Move) String() string {
	return fmt.Sprintf("%v -> %v Piece:%q Promoted:%q Capture:%t Double:%t Enpassant:%t Castling:%t",
		m.getSource(), m.getTarget(), pieceToChar[m.getPiece()], pieceToChar[m.getPromotionPiece()],
		m.isCapture(), m.isDoublePawnPush(), m.isEnpassant(), m.isCastling(),
	)
}

func TestGenerateMoves(t *testing.T) {
	Init()

	board := NewBoadFromFen([]byte(tricky_position))

	moves := Moves{}

	board.generateMoves(&moves)

	if moves.count != 48 {
		t.Fatalf("Expected 48 moves but got %d\n", moves.count)
	}
}
