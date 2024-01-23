package application

import "net/http"

func getUserID(r *http.Request) (string, error) {
	cookieUser, err := r.Cookie("user_id")
	if err != nil {
		return "", err
	}

	return cookieUser.Value, nil
}
