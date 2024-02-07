package application

import (
	"log"
	"net/http"

	"github.com/LyubenGeorgiev/shah/util"
	"github.com/LyubenGeorgiev/shah/view/account"

	"github.com/LyubenGeorgiev/shah/db"

)

func (app *App) HandleProfile(w http.ResponseWriter, r *http.Request) {

	userID, err := util.GetUserID(r)
	if err != nil || userID == "" {
		log.Println("Unknown user!", userID, "Error", err.Error())
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	user, err :=  db.NewPostgresStorage().FindByUserID(userID)

	account.Account(user).Render(r.Context(), w)

}
