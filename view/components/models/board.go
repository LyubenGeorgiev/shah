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
	Unselect    []int
	View        Side
}

// type Action string

// const (
// 	SELECT   Action = "select"
// 	UNSELECT Action = "unselect"
// 	MOVE     Action = "move"
// 	CAPTURE  Action = "capture"
// 	NONE     Action = ""
// )
