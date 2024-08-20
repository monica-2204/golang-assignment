package utils

import (
	"context"
	"net/http"
	"strings"
)

// This is used to assign values to created_by and updated_by in transport/student
func GetCurrentUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("userID").(string); ok {
		return userID
	}
	return ""
}

// This is used extract token to in userIDMiddleware in transport/middleware
func ExtractTokenFromHeader(r *http.Request) string {

	authHeader := r.Header.Get("Authorization")

	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}
