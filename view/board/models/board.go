package models

type Side int

const (
	White Side = iota
	Black
)

type BoardState struct {
	Highlighted []int
	Pieces      map[int]string
	Moves       []int
	Captures    []int
	View        Side
}

type Message struct {
	Text           string
	IsFromOpponent bool
}
