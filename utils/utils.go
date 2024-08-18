package utils

import (
	"context"
	"net/http"
	"strings"
)

// getCurrentUserID extracts the user ID from the context
func GetCurrentUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("userID").(string); ok {
		return userID
	}
	return ""
}

// extractTokenFromHeader extracts the JWT token from the Authorization header.
func ExtractTokenFromHeader(r *http.Request) string {

	authHeader := r.Header.Get("Authorization")

	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}
