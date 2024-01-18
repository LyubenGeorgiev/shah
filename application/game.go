package application

import (
	"fmt"
	"net/http"

	"github.com/LyubenGeorgiev/shah/view/components/models"
	"github.com/LyubenGeorgiev/shah/view/layout"
)

// TODO make actual handler :D
func (a *App) Game(w http.ResponseWriter, r *http.Request) {
	err := layout.BoardTest(&models.BoardState{
		Pieces: [...]string{
			"r", "n", "b", "q", "k", "b", "n", "r",
			"p", "p", "p", "p", "p", "p", "p", "p",
			"empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty",
			"empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty",
			"empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty",
			"empty", "empty", "empty", "empty", "empty", "empty", "empty", "empty",
			"P", "P", "P", "P", "P", "P", "P", "P",
			"R", "N", "B", "Q", "K", "B", "N", "R",
		},
		Clicable: [...]bool{
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			true, true, true, true, true, true, true, true,
			true, true, true, true, true, true, true, true,
		},
		View: models.White,
		Side: models.White,
	}).Render(r.Context(), w)
	if err != nil {
		fmt.Printf("Error rendering game board: %v", err)
	}
}
