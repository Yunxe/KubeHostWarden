package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func Cors(next http.Handler) http.Handler {
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}

	return cors.New(corsOptions).Handler(next)
}
