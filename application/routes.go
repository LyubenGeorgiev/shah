package application

import (
	"fmt"
	"net/http"

	"github.com/LyubenGeorgiev/shah/view/layout"
)

func (app *App) loadRoutes() {
	app.router.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))
	app.router.PathPrefix("/static/images/").Handler(http.StripPrefix("/static/images/", http.FileServer(http.Dir("static/images"))))

	app.router.Use(app.authMiddleware)

	app.router.HandleFunc("/register", app.RegistrationFrom).Methods("GET")
	app.router.HandleFunc("/register", app.Register).Methods("POST")

	app.router.HandleFunc("/login", app.LoginFrom).Methods("GET")
	app.router.HandleFunc("/login", app.Login).Methods("POST")

	app.router.Handle("/logout", app.requiredAuthMiddleware(http.HandlerFunc(app.Logout))).Methods("GET")

	app.router.HandleFunc("/play", app.Manager.HandlePlay).Methods("GET")
	app.router.HandleFunc("/computer", app.Computer).Methods("GET")
	app.router.HandleFunc("/computer/{gameID}", app.ComputerGame).Methods("GET")
	app.router.HandleFunc("/tournaments", app.Tournaments).Methods("GET")
	app.router.HandleFunc("/game/{id}", app.Manager.HandleGame).Methods("GET")

	app.router.HandleFunc("/news", app.News).Methods("GET")
	app.router.Handle("/createNews", app.requiredAdminRole(http.HandlerFunc(app.CreateNews))).Methods("GET")
	app.router.Handle("/createNews", app.requiredAdminRole(http.HandlerFunc(app.NewNews))).Methods("POST")

	app.router.HandleFunc("/replays/{gameID}", app.HandleReplays).Methods("GET")
	app.router.HandleFunc("/replays/{gameID}/{move}", app.HandleReplaysMove).Methods("GET")

	app.router.HandleFunc("/history/{userID}/{page}", app.HandleMatchHistory).Methods("GET")

	app.router.HandleFunc("/messages", app.HandleChatsLayout).Methods("GET")
	app.router.HandleFunc("/loadchats/{page}", app.HandleLoadChats).Methods("GET")
	app.router.HandleFunc("/messages/{userID}/{page}", app.HandleMessages).Methods("GET")
	app.router.HandleFunc("/chats/{userID}", app.HandleChats).Methods("GET")
	app.router.HandleFunc("/chats/{userID}", app.HandleChatsWrite).Methods("POST")

	app.router.HandleFunc("/search", app.HandleSearch).Methods("GET")
	app.router.HandleFunc("/profilewidgets/{id}", app.HandleProfilewidgets).Methods("GET")

	app.router.Handle("/account", app.requiredAuthMiddleware(http.HandlerFunc(app.HandleAccount))).Methods("GET")
	app.router.HandleFunc("/profiles/{id}", app.HandleProfiles).Methods("GET")
	app.router.Handle("/upload", app.requiredAuthMiddleware(http.HandlerFunc(app.HandleUpload))).Methods("PUT")
	app.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := layout.Home().Render(r.Context(), w)
		if err != nil {
			fmt.Printf("Error rendering at home page: %v", err)
		}
	}).Methods("GET")
}
