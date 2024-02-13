package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/LyubenGeorgiev/shah/db/models"
	"github.com/LyubenGeorgiev/shah/util"
	"github.com/LyubenGeorgiev/shah/view/login"
	"github.com/LyubenGeorgiev/shah/view/registration"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

const (
	auth_duration      = 24 * time.Hour
	half_auth_duration = 12 * time.Hour
)

func (a *App) RegistrationFrom(w http.ResponseWriter, r *http.Request) {
	err := registration.Register().Render(r.Context(), w)
	if err != nil {
		fmt.Printf("Error rendering at registration form page: %v", err)
	}
}

func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{Rating: 1000, Image: "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b5/Windows_10_Default_Profile_Picture.svg/2048px-Windows_10_Default_Profile_Picture.svg.png", Role: "USER"}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Decoding body failed during registration", http.StatusInternalServerError)

		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Password Encryption failed", http.StatusInternalServerError)

		return
	}

	user.Password = string(pass)

	if err := a.Storage.CreateUser(user); err != nil {
		http.Error(w, "Saving registration failed", http.StatusInternalServerError)

		return
	}

	w.Header().Set("HX-Location", "/login")
}

func (a *App) LoginFrom(w http.ResponseWriter, r *http.Request) {
	err := login.Login().Render(r.Context(), w)
	if err != nil {
		fmt.Printf("Error rendering at login form page: %v", err)
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

	id, err := a.Storage.FindOneUser(user.Email, user.Password)
	if err != nil {
		log.Println("Error during login:", err)
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	// Generate a unique token (in this example, using UUID)
	token := uuid.New().String()
	redisKey := fmt.Sprintf("auth:%d", id)

	expiration := time.Now().Add(auth_duration) // Adjust the expiration time as needed

	err = a.Cache.Set(r.Context(), redisKey, token, auth_duration)
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

		// Extend user authentication
		if cookieAuth.Expires.Sub(time.Now()) < half_auth_duration {
			err = a.Cache.Set(r.Context(), redisKey, authToken, auth_duration)
			if err != nil {
				http.Error(w, "Error storing auth token", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    authToken,
				Expires:  time.Now().Add(auth_duration),
				HttpOnly: true,
				Path:     "/",
			})
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

func (a *App) requiredAdminRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := util.GetUserID(r)
		if err != nil {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		user, err := a.Storage.FindByUserID(userID)
		if err != nil || user.Role != "ADMIN" {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserID(r)
	if err != nil || userID == "" {
		http.Error(w, "Unknown user!", http.StatusUnauthorized)
		return
	}

	redisKey := fmt.Sprintf("auth:%s", userID)

	err = a.Cache.Del(r.Context(), redisKey)
	if err != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
