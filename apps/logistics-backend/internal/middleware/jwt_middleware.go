package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	ContextUserID contextKey = "userID"
	ContextRole   contextKey = "role"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			http.Error(w, "Server misconfigured (no JWT_SECRET)", http.StatusInternalServerError)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Check the signing method is HMAC (HS256)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Extract user ID and role from claims
		userID, ok1 := claims["sub"].(string)
		role, ok2 := claims["role"].(string)
		if !ok1 || !ok2 {
			http.Error(w, "Missing token claims", http.StatusUnauthorized)
			return
		}

		// Add to context
		ctx := context.WithValue(r.Context(), ContextUserID, userID)
		ctx = context.WithValue(ctx, ContextRole, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAdminIDFromContext(ctx context.Context) (uuid.UUID, error) {
	role, ok := ctx.Value(ContextRole).(string)
	if !ok || role != "admin" {
		return uuid.Nil, errors.New("not authorized as admin")
	}

	idStr, ok := ctx.Value(ContextUserID).(string)
	if !ok {
		return uuid.Nil, errors.New("missing user ID in context")
	}

	return uuid.Parse(idStr)
}

func GetDriverIDFromContext(ctx context.Context) (uuid.UUID, error) {
	role, ok := ctx.Value(ContextRole).(string)
	if !ok || role != "driver" {
		return uuid.Nil, errors.New("not authorized as driver")
	}

	idStr, ok := ctx.Value(ContextUserID).(string)
	if !ok {
		return uuid.Nil, errors.New("missing user ID in context")
	}

	return uuid.Parse(idStr)
}
