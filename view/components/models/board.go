package models

type Side int

const (
	White Side = iota
	Black
)

type BoardState struct {
	Pieces [64]string
	View   Side
	Side   Side
}
