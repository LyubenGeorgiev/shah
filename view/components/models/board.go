package models

type Side int

const (
	White Side = iota
	Black
)

type BoardState struct {
	Pieces   [64]string
	Clicable [64]bool
	View     Side
	Side     Side
}
