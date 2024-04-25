package user

import (
	"context"
	"kubehostwarden/db"
	"kubehostwarden/types"
	"kubehostwarden/utils/constant"
	resp "kubehostwarden/utils/responsor"
	"net/http"
	"net/url"
)

type retrieveResp struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func Retrieve(ctx context.Context, values url.Values) resp.Responsor {
	userId := ctx.Value(constant.UserIDKey).(string)
	var user types.User
	res := db.GetMysqlClient().Client.WithContext(ctx).Where("id = ?", userId).First(&user)
	if res.Error != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: "failed to retrieve user",
		}
	}

	response := retrieveResp{
		Username: user.Username,
		Email:    user.Email,
	}

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "user retrieved successfully",
		Result:  response,
	}
}
