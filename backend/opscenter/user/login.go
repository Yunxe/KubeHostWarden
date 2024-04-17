package user

import (
	"context"
	"fmt"
	"kubehostwarden/db"
	"kubehostwarden/types"
	resp "kubehostwarden/utils/responsor"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResp struct {
	Token string `json:"token"`
}

func Login(ctx context.Context, loginReq LoginReq) resp.Responsor {
	var existedUser *types.User
	if result := db.GetMysqlClient().Client.WithContext(ctx).Where("email = ?", loginReq.Email).First(&existedUser); result.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to query user: %v", result.Error),
		}
	}

	if existedUser == nil {
		return resp.Responsor{
			Code:    http.StatusNotFound,
			Message: "user not found",
		}
	}

	if existedUser.Password != loginReq.Password {
		return resp.Responsor{
			Code:    http.StatusUnauthorized,
			Message: "wrong password",
		}
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    existedUser.Id,
		"email": loginReq.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		// Handle the error here
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("failed to sign token: %v", err),
		}
	}

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "login successfully",
		Result: LoginResp{
			Token: tokenString,
		},
	}
}
