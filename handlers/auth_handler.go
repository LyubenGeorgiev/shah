package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LyubenGeorgiev/shah/db"
	"github.com/LyubenGeorgiev/shah/db/models"
	"github.com/LyubenGeorgiev/shah/view/registration"

	"golang.org/x/crypto/bcrypt"
)

type ErrorResponse struct {
	Err string
}

type error interface {
	Error() string
}

type AuthHandler struct {
	Storage db.Storage
}

func NewAuthHandler(storage db.Storage) *AuthHandler {
	return &AuthHandler{
		Storage: storage,
	}
}

func (h *AuthHandler) RegistrationFrom(w http.ResponseWriter, r *http.Request) {
	registration.Register().Render(r.Context(), w)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		fmt.Println(err)
		err := ErrorResponse{
			Err: "Decoding body failed during registration",
		}
		json.NewEncoder(w).Encode(err)

		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		err := ErrorResponse{
			Err: "Password Encryption failed",
		}
		json.NewEncoder(w).Encode(err)

		return
	}

	user.Password = string(pass)

	if err := h.Storage.CreateUser(user); err != nil {
		fmt.Println(err)
		err := ErrorResponse{
			Err: "Storing registration in database failed",
		}
		json.NewEncoder(w).Encode(err)

		return
	}
}
