package logger

import (
	"context"
	"kubehostwarden/utils/constant"
	resp "kubehostwarden/utils/responsor"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func GetLogs(ctx context.Context, values url.Values) resp.Responsor {
	userId := ctx.Value(constant.UserIDKey).(string)

	logFilePath := filepath.Join("logs", userId+".log")
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		return resp.Responsor{
			Code:    http.StatusNotFound,
			Message: "log file not found",
		}
	}

	logData, err := os.ReadFile(logFilePath)
	if err != nil {
		return resp.Responsor{
			Code:    http.StatusInternalServerError,
			Message: "failed to read log file",
		}
	}

	return resp.Responsor{
		Code:    http.StatusOK,
		Message: "log file fetched successfully",
		Result:  string(logData),
	}
}
