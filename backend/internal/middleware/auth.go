package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/razvan/library-app/internal/utils"
	"github.com/razvan/library-app/pkg/auth"
)

type contextKey string

const (
	UserIDKey  contextKey = "user_id"
	EmailKey   contextKey = "email"
	IsAdminKey contextKey = "is_admin"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		token := parts[1]
		claims, err := auth.ValidateToken(token, m.jwtSecret)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)
		ctx = context.WithValue(ctx, IsAdminKey, claims.IsAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(IsAdminKey).(bool)
		if !ok || !isAdmin {
			utils.ErrorResponse(w, http.StatusForbidden, "admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}

func GetEmail(ctx context.Context) string {
	email, ok := ctx.Value(EmailKey).(string)
	if !ok {
		return ""
	}
	return email
}

func IsAdmin(ctx context.Context) bool {
	isAdmin, ok := ctx.Value(IsAdminKey).(bool)
	if !ok {
		return false
	}
	return isAdmin
}
