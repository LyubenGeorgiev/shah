package application

import (
	"net/http"

	"github.com/LyubenGeorgiev/shah/view/layout"
)

func (a *App) Tournaments(w http.ResponseWriter, r *http.Request) {
	// TODO implement this

	layout.Tournaments().Render(r.Context(), w)
}
