package game

import (
	"github.com/LyubenGeorgiev/shah/chess"
)

type Game struct {
	GameID  string
	Board   chess.Board
	actions chan<- string
	white   Client
	black   Client
}
