package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// implement Auth function
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the token from the request header
		token := r.Header.Get("Authorization")
		if token == "" {
			// if the token is empty, return 401 Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// validate the token
		id, email, err := validateToken(token); 
		if err != nil {
			// if the token is invalid, return 401 Unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("X-User-Id", id)
		w.Header().Set("X-User-Email", email)
		// if the token is valid, call the next handler
		next.ServeHTTP(w, r)
	})
}

// validateToken is a dummy function to validate the token
func validateToken(tokenString string) (string, string, error) {
	// Check if the token has the "Bearer " prefix
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return "", "", fmt.Errorf("invalid token")
	}

	// Remove "Bearer " from the token string
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Define a new JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token's signature algorithm is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for token verification
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", "", err
	}

	// Check if the token claims are valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["id"].(string), claims["email"].(string), nil
	} else {
		return "", "", fmt.Errorf("invalid token")
	}
}
