package chess

import (
	"testing"
)

var nodes uint64 = 0
var board *Board

// perft driver
func perftDriver(depth int) {
	// create move list instance
	moves := &Moves{}

	// generate moves
	board.generateMoves(moves)

	// loop over generated moves
	for i := 0; i < moves.count; i++ {
		// preserve board state
		copy := *board

		// make move
		if !board.MakeMove(moves.moves[i], false) {
			// skip to the next move
			continue
		}

		// call perft driver recursively
		if depth != 1 {
			perftDriver(depth - 1)
		} else {
			nodes++
		}

		// take back
		*board = copy
	}
}

func TestPerft(t *testing.T) {
	board = NewBoadFromFen([]byte(start_position))
	nodes = 0

	// perft
	perftDriver(6)

	if nodes != 119060324 {
		t.Log("Expected 119060324 nodes for initial board at depth 6 but got:", nodes)
		t.Fail()
	}

	board = NewBoadFromFen([]byte(tricky_position))
	nodes = 0

	// perft
	perftDriver(5)

	if nodes != 193690690 {
		t.Log("Expected 193690690 nodes for tricky board at depth 5 but got:", nodes)
		t.Fail()
	}
}
