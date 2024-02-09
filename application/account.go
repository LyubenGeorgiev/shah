package application

import (
	"fmt"
	"net/http"

	"github.com/LyubenGeorgiev/shah/util"
	"github.com/LyubenGeorgiev/shah/view/account"
	"github.com/gorilla/mux"

	"github.com/LyubenGeorgiev/shah/db"

	"encoding/base64"
	"io"
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
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/profiles/"+userID, http.StatusSeeOther)
}





// HandleUpload handles the upload of an image for a user account
func (app *App) HandleUpload(w http.ResponseWriter, r *http.Request) {

	fmt.Println("vliza")
    // Parse the form data
    err := r.ParseMultipartForm(10 << 20) // 10 MB limit
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Get the user ID from the request
    userID, err := util.GetUserID(r)
    if err != nil || userID == "" {
        http.Error(w, "Unknown user!", http.StatusUnauthorized)
        return
    }

    // Get the uploaded file from the form
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Unable to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Read the contents of the file
    fileBytes, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "Unable to read file", http.StatusInternalServerError)
        return
    }

    // Encode the file contents to base64
    encodedFile := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(fileBytes)

    // Update the user's image in the database
    if err := db.NewPostgresStorage().UpdateUserImage(userID, encodedFile); err != nil {
        http.Error(w, "Failed to update user image", http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "<script>window.location.href='/profiles/%s';</script>", userID)

}
