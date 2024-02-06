package util

import "net/http"

func GetUserID(r *http.Request) (string, error) {
	cookieUser, err := r.Cookie("user_id")
	if err != nil {
		return "", err
	}

	return cookieUser.Value, nil
}
