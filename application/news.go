package application

import (
	"encoding/base64" // For base64 encoding and decoding
	"io" // For reading file content
	"net/http"

	"github.com/LyubenGeorgiev/shah/db/models"
	"github.com/LyubenGeorgiev/shah/view/news"
	"github.com/LyubenGeorgiev/shah/util"
)

func (a *App) News(w http.ResponseWriter, r *http.Request) {

	userID, err := util.GetUserID(r)

	if err != nil || userID == "" {
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	user, err := a.Storage.FindByUserID(userID)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}


	isAdmin := user.Role == "ADMIN"

	if err != nil {
		http.Error(w, "Unknown user!", http.StatusNotFound)
		return
	}

	newsList, err := a.Storage.GetAllNews()
	if err != nil {
		http.Error(w, "Failed to retrieve news items", http.StatusInternalServerError)
		return
	}
	news.News(newsList , isAdmin).Render(r.Context(), w)
}


func (a *App) CreateNews(w http.ResponseWriter, r *http.Request) {

	news.CreateNews().Render(r.Context(), w)
}


func (app *App) NewNews(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	title := r.FormValue("title")
	description := r.FormValue("description")
	url := r.FormValue("url")

	// Get the uploaded file from the form
	file, _, err := r.FormFile("image")
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

	// Create a new News object
	news := &models.News{
		Title:       title,
		Description: description,
		URL:         url,
		Image:       encodedFile,
	}

	// Update the user's image in the database
	if err := app.Storage.CreateNews(news); err != nil {
		http.Error(w, "Failed to update user image", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("HX-Location", "/news")
}
