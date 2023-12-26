package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LyubenGeorgiev/shah/db"
	"github.com/LyubenGeorgiev/shah/handlers"
	"github.com/gorilla/mux"
)

// swagger:response HelloResponse
type HelloResponse struct {
	// in:body
	Message string `json:"message"`
}

// swagger:route GET / hello
// helloHandler responds with a "Hello, World!" message.
// Consumes:
//   - application/json
//
// Produces:
//   - application/json
//
// Responses:
//
//	200: HelloResponse
func helloHandler(rw http.ResponseWriter, r *http.Request) {
	response := HelloResponse{Message: "Hello, World!"}
	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		http.Error(rw, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	ah := handlers.NewAuthHandler(db.NewPostgresStorage())

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/register", ah.RegistrationFrom).Methods("GET")
	r.HandleFunc("/register", ah.Register).Methods("POST")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "HOME PAGE")
	}).Methods("GET")

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
