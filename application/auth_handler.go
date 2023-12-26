package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func (a *App) RegistrationFrom(w http.ResponseWriter, r *http.Request) {
	registration.Register().Render(r.Context(), w)
}

func (a *App) Register(w http.ResponseWriter, r *http.Request) {
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

	if err := a.Storage.CreateUser(user); err != nil {
		fmt.Println(err)
		err := ErrorResponse{
			Err: "Storing registration in database failed",
		}
		json.NewEncoder(w).Encode(err)

		return
	}
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println("Decoding body failed during login:", err)
		http.Error(w, "Decoding body failed during login", http.StatusBadRequest)

		return
	}

	if err = a.Storage.FindOneUser(user.Email, user.Password); err != nil {
		log.Println("Error during login:", err)
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

}
