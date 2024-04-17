package user

import (
	"encoding/json"
	"kubehostwarden/db"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type registerInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var regInfo registerInfo
	err := json.NewDecoder(r.Body).Decode(&regInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Check if the email already exists
	var existingUser User
	db.GetMysqlClient().Client.WithContext(r.Context()).Where("email = ?", regInfo.Email).First(&existingUser)
	if existingUser.Email != "" {
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(map[string]string{"error": "email already registered"})
		return
	}

	user := &User{
		Id:       uuid.NewString()[:8],
		Username: regInfo.Username,
		Password: regInfo.Password,
		Email:    regInfo.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the new user
	db.GetMysqlClient().Client.Save(&user)

	json.NewEncoder(w).Encode(map[string]string{"message": "successfully registered"})
}
