package utils

import "context"

// getCurrentUserID extracts the user ID from the context
func GetCurrentUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("userID").(string); ok {
		return userID
	}
	return ""
}
