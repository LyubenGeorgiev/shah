package application

import (
	"net/http"

	"github.com/LyubenGeorgiev/shah/view/nav"
)

func (a *App) HandleSearch(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		return
	}

	users, _ := a.Storage.FetchUsersByUsername(keyword)

	nav.SearchResults(users).Render(r.Context(), w)
}
