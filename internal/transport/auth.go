package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

// JWTAuth - a handy middleware function that will provide basic auth around specific endpoints
func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorizedResponse(w, "Authorization header missing")
			log.Error("Authorization header missing")
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			unauthorizedResponse(w, "Invalid authorization header format")
			log.Error("Invalid authorization header format")
			return
		}

		tokenStr := authHeaderParts[1]
		if validateToken(tokenStr) {
			next(w, r)
		} else {
			unauthorizedResponse(w, "Invalid JWT token")
			log.Error("Invalid JWT token")
			return
		}
	}
}

// validateToken - validates an incoming JWT token
func validateToken(accessToken string) bool {
	mySigningKey := []byte(os.Getenv("JWT_SECRET"))
	if len(mySigningKey) == 0 {
		log.Error("JWT_SECRET environment variable not set")
		return false
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return mySigningKey, nil
	})

	if err != nil {
		log.Errorf("Error parsing JWT token: %v", err)
		return false
	}

	return token.Valid
}

// unauthorizedResponse - sends a JSON unauthorized response
func unauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
