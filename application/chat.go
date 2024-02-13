package application

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/LyubenGeorgiev/shah/db/models"
	"github.com/LyubenGeorgiev/shah/util"
	"github.com/LyubenGeorgiev/shah/view/messages"
	"github.com/gorilla/mux"
)

func (app *App) HandleLoadChats(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserID(r)
	if err != nil || userID == "" {
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	page, err := strconv.Atoi(mux.Vars(r)["page"])
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	chatsIDs, err := app.Storage.GetRecentChatsUserIDs(userID, page, 10)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	messages.ShowChats(page, chatsIDs).Render(r.Context(), w)
}

func (app *App) HandleMessages(w http.ResponseWriter, r *http.Request) {
	userID1, err := util.GetUserID(r)
	if err != nil || userID1 == "" {
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	userID2 := mux.Vars(r)["userID"]
	page, err := strconv.Atoi(mux.Vars(r)["page"])
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	msgs, err := app.Storage.GetRecentMessagesWith(userID1, userID2, page, 10)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	messages.ShowMessages(page, userID2, msgs).Render(r.Context(), w)
}

func (app *App) HandleChats(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]

	messages.Chat(userID).Render(r.Context(), w)
}

func (app *App) HandleChatsWrite(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserID(r)
	if err != nil || userID == "" {
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	from, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	toStr := mux.Vars(r)["userID"]
	to, err := strconv.ParseUint(toStr, 10, 64)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	bytes, _ := io.ReadAll(r.Body)

	msg := models.Message{From: uint(from), To: uint(to)}
	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = app.Storage.CreateMessage(&msg)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	messages.Chat(toStr).Render(r.Context(), w)
}

func (app *App) HandleChatsLayout(w http.ResponseWriter, r *http.Request) {
	messages.Layout().Render(r.Context(), w)
}
