package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/LyubenGeorgiev/shah/db/models"
	"github.com/LyubenGeorgiev/shah/view/login"
	"github.com/LyubenGeorgiev/shah/view/registration"
	"github.com/google/uuid"

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

	w.Header().Set("HX-Redirect", "/login")
}

func (a *App) LoginFrom(w http.ResponseWriter, r *http.Request) {
	login.Login().Render(r.Context(), w)
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println("Decoding body failed during login:", err)
		http.Error(w, "Decoding body failed during login", http.StatusBadRequest)

		return
	}

	id, err := a.Storage.FindOneUser(user.Email, user.Password)
	if err != nil {
		log.Println("Error during login:", err)
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	// Generate a unique token (in this example, using UUID)
	token := uuid.New().String()
	redisKey := fmt.Sprintf("auth:%d", id)

	duration := 24 * time.Hour
	expiration := time.Now().Add(duration) // Adjust the expiration time as needed

	fmt.Println("Set", redisKey, token)
	err = a.Cache.Set(r.Context(), redisKey, token, duration)
	if err != nil {
		http.Error(w, "Error storing auth token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprint(id),
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", "/")
}

func (a *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the cookie named "auth_token"
		cookieAuth, err := r.Cookie("auth_token")
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "authenticated", false)))
			return
		}
		cookieUser, err := r.Cookie("user_id")
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "authenticated", false)))
			return
		}

		// Get the token value from the cookie
		authToken := cookieAuth.Value
		userID := cookieUser.Value
		redisKey := fmt.Sprintf("auth:%s", userID)

		if !a.Cache.Exists(r.Context(), redisKey, authToken) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "authenticated", false)))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "authenticated", true)))
	})
}

func (a *App) requiredAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value("authenticated").(bool) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		}
	})
}
