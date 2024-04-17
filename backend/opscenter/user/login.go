package user

import (
	"encoding/json"
	"kubehostwarden/db"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type loginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginInfo loginInfo
	err := json.NewDecoder(r.Body).Decode(&loginInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	var existedUser *User
	db.GetMysqlClient().Client.WithContext(r.Context()).Where("email = ?", loginInfo.Email).First(&existedUser)

	if existedUser == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "this email has not been registered!"})
		return
	}

	if existedUser.Password != loginInfo.Password {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "wrong password!"})
		return
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": loginInfo.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("your_secret_key")) // Replace "your_secret_key" with your secret key
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate token"})
		return
	}

	json.NewEncoder(w).Encode(loginResponse{Token: tokenString})
}
