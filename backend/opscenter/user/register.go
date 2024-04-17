package user

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/types"
	resp "kubehostwarden/utils/responsor"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type registerReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func Register(ctx context.Context, regInfo registerReq) resp.Responsor {
	// Check if the email already exists
	var existingUser types.User
	db.GetMysqlClient().Client.WithContext(ctx).Where("email = ?", regInfo.Email).First(&existingUser)
	if existingUser.Email != "" {
		return resp.Responsor{
			Code:    http.StatusConflict,
			Message: "email already exists",
		}
	}

	user := &types.User{
		Id:        uuid.NewString(),
		Username:  regInfo.Username,
		Password:  regInfo.Password,
		Email:     regInfo.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save the new user
	res := db.GetMysqlClient().Client.Save(&user)
	if res.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to save user: %v", res.Error),
		}
	}

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "user registered successfully",
		Result:  user,
	}
}
