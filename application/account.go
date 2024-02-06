package application

import (
	"net/http"

	"github.com/LyubenGeorgiev/shah/view/account"
)

func (app *App)Account(w http.ResponseWriter, r *http.Request) {

	account.Account().Render(r.Context(), w)

}
