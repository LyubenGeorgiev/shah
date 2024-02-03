package chess

type Engine struct {
	Board
	TimeController
	Searcher
	RepetitionTable []Bitboard
	ply             int
	nodes           uint64
}

func NewEngine() Engine {
	return Engine{
		Board:           Board{},
		TimeController:  NewTimeController(),
		RepetitionTable: []Bitboard{},
	}
}
