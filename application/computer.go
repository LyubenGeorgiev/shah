package application

import (
	"net/http"

	"github.com/LyubenGeorgiev/shah/view/layout"
)

func (a *App) Computer(w http.ResponseWriter, r *http.Request) {
	// TODO implement this

	layout.Computer().Render(r.Context(), w)
}
