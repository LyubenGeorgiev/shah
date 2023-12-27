package application

import (
	"net/http"

	"github.com/LyubenGeorgiev/shah/view/layout"
)

func (app *App) loadRoutes() {
	app.router.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))
	app.router.PathPrefix("/static/images/").Handler(http.StripPrefix("/static/images/", http.FileServer(http.Dir("static/images"))))

	app.router.HandleFunc("/register", app.RegistrationFrom).Methods("GET")
	app.router.HandleFunc("/register", app.Register).Methods("POST")

	app.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		layout.Home().Render(r.Context(), w)
	}).Methods("GET")
}
