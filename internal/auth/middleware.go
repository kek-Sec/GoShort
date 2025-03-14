package auth

import (
	"GoShort/pkg/logger"
	"net/http"
)

// AuthMiddleware checks if the request has a valid JWT token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the request
		tokenString, err := ExtractTokenFromRequest(r)
		if err != nil {
			logger.Error("Authentication failed: "+err.Error(), nil)
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Verify the token
		claims, err := VerifyToken(tokenString)
		if err != nil {
			logger.Error("Token verification failed: "+err.Error(), nil)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to the request context
		ctx := r.Context()
		ctx = AddClaimsToContext(ctx, claims)

		// Call the next handler with the updated context
		next(w, r.WithContext(ctx))
	}
}

// AdminOnlyMiddleware ensures that only users with admin role can access the endpoint
func AdminOnlyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		claims, err := GetClaimsFromContext(r.Context())
		if err != nil || claims.Role != "admin" {
			logger.Error("Admin access denied", map[string]interface{}{
				"username": claims.Username,
				"role":     claims.Role,
			})
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		next(w, r)
	})
}
