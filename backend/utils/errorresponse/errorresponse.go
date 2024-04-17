package errorresponse

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteErrorResponse(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(http.StatusOK)

	errorResponse := ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}

	if writeErr := json.NewEncoder(w).Encode(errorResponse); writeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
