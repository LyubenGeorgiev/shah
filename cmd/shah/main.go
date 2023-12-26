package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LyubenGeorgiev/shah/db"
	"github.com/LyubenGeorgiev/shah/handlers"
	"github.com/LyubenGeorgiev/shah/view/layout"
	"github.com/gorilla/mux"
)

func main() {
	ah := handlers.NewAuthHandler(db.NewPostgresStorage())

	r := mux.NewRouter().StrictSlash(true)
	r.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))
	r.PathPrefix("/static/images/").Handler(http.StripPrefix("/static/images/", http.FileServer(http.Dir("static/images"))))

	r.HandleFunc("/register", ah.RegistrationFrom).Methods("GET")
	r.HandleFunc("/register", ah.Register).Methods("POST")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		layout.Base("Shah.com - Play Chess Online").Render(r.Context(), w)
	}).Methods("GET")

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
