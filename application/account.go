package application

import (
	"log"
	"net/http"

	"github.com/LyubenGeorgiev/shah/util"
	"github.com/LyubenGeorgiev/shah/view/account"
	"github.com/gorilla/mux"

	"github.com/LyubenGeorgiev/shah/db"
)

func (app *App) HandleProfiles(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	user, err := db.NewPostgresStorage().FindByUserID(userID)
	if err != nil {
		http.Error(w, "Unknown user!", http.StatusNotFound)
		return
	}

	account.Account(user).Render(r.Context(), w)
}

func (app *App) HandleAccount(w http.ResponseWriter, r *http.Request) {

	userID, err := util.GetUserID(r)
	if err != nil || userID == "" {
		log.Println("Unknown user!", userID, "Error", err.Error())
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "profiles/"+userID, http.StatusSeeOther)
}
