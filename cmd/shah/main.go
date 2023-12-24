package shah

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LyubenGeorgiev/shah/db"
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
	db.SetupDatabase()
	http.HandleFunc("/", helloHandler)

	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
