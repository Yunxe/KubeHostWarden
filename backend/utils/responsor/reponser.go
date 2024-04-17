package responsor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kubehostwarden/utils/middleware"
	"net/http"
	"net/url"
)

type Responsor struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func Decode[T any](r io.Reader) (*T, error) {
	var result T
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func HandleGet(handler func(context.Context, url.Values) Responsor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		query := r.URL.Query()
		resp := handler(r.Context(), query)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			resp := Responsor{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to encode response: %v", err),
				Result:  nil,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
}

func HandlePost[reqType any](handlerFunc func(ctx context.Context, req reqType) Responsor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		req, err := Decode[reqType](r.Body)
		if err != nil {
			resp := Responsor{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("invalid request body: %v", err),
				Result:  nil,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		err = middleware.ValidateReq(*req)
		if err != nil {
			resp := Responsor{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("invalid request body: %v", err),
				Result:  nil,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := handlerFunc(r.Context(), *req)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			resp := Responsor{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("failed to encode response: %v", err),
				Result:  nil,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
}
