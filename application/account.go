package application

import (
	"net/http"
	"strconv"

	"github.com/LyubenGeorgiev/shah/util"
	"github.com/LyubenGeorgiev/shah/view/account"
	components "github.com/LyubenGeorgiev/shah/view/board"
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

	curUserID, _ := util.GetUserID(r)

	account.Account(user, userID == curUserID).Render(r.Context(), w)
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
	if err := app.Storage.UpdateUserImage(userID, encodedFile); err != nil {
		http.Error(w, "Failed to update user image", http.StatusInternalServerError)
		return
	}

	account.ProfilePicture(encodedFile).Render(r.Context(), w)
}

// profilewidgets
func (app *App) HandleProfilewidgets(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	user, err := app.Storage.FindByUserID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	components.Profilewidget(user).Render(r.Context(), w)
}

func (app *App) HandleMatchHistory(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]
	page, err := strconv.Atoi(mux.Vars(r)["page"])
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	games, err := app.Storage.GetMatchHistoryGames(userID, page, 10)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	account.ShowHistory(userID, page, games).Render(r.Context(), w)
}


func (app *App) HandleUsers(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(mux.Vars(r)["page"])
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	users, err := app.Storage.GetAllUsers(page, 10)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	account.ShowUsers( page, users).Render(r.Context(), w)
}

func(app  *App) HandleUsersShow(w http.ResponseWriter, r *http.Request){
	account.UsersPage().Render(r.Context(),w)
}

func(app  *App) HandleDeleteUser(w http.ResponseWriter, r *http.Request){
	id, err := strconv.ParseUint(mux.Vars(r)["id"] , 10, 64)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = app.Storage.DeleteUserByID(uint(id))

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Location", "/users")



}


func(app  *App) HandleUpdateRole(w http.ResponseWriter, r *http.Request){
	id, err := strconv.ParseUint(mux.Vars(r)["id"] , 10, 64)
	
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	role := mux.Vars(r)["role"]


	err = app.Storage.UpdateUser(uint(id) ,role)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Location", "/users")



}

