package application

import (
	"fmt"
	"net/http"

	"github.com/LyubenGeorgiev/shah/view/layout"
)

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a protected endpoint.")
}

func (app *App) loadRoutes() {
	app.router.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))
	app.router.PathPrefix("/static/images/").Handler(http.StripPrefix("/static/images/", http.FileServer(http.Dir("static/images"))))

	app.router.Use(app.authMiddleware)

	app.router.HandleFunc("/register", app.RegistrationFrom).Methods("GET")
	app.router.HandleFunc("/register", app.Register).Methods("POST")

	app.router.HandleFunc("/login", app.LoginFrom).Methods("GET")
	app.router.HandleFunc("/login", app.Login).Methods("POST")

	app.router.Handle("/protected", app.requiredAuthMiddleware(http.HandlerFunc(protectedHandler)))

	app.router.HandleFunc("/play", app.Manager.HandlePlay).Methods("GET")
	app.router.HandleFunc("/computer", app.Computer).Methods("GET")
	app.router.HandleFunc("/tournaments", app.Tournaments).Methods("GET")
	app.router.HandleFunc("/game/{id}", app.Manager.HandleGame).Methods("GET")

	app.router.HandleFunc("/account", app.Account ).Methods("GET")

	app.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := layout.Home().Render(r.Context(), w)
		if err != nil {
			fmt.Printf("Error rendering at home page: %v", err)
		}
	}).Methods("GET")
}
