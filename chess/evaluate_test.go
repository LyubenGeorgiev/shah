package chess

import "testing"

func TestEvaluate(t *testing.T) {
	Init()

	board := NewBoadFromFen([]byte(start_position))

	if board.evaluate() != 0 {
		t.Log("Expected evaluation 0 for start position but got:", board.evaluate())
		t.Fail()
	}
}
