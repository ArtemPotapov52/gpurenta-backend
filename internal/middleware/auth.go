package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ArtemPotapov52/gpurenta/internal/auth"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := auth.ValidateToken(tokenStr, jwtSecret)
			if err != nil {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			ctx := NewContextWithUser(r.Context(), claims.UserID, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AgentAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agentID := r.Header.Get("X-Agent-ID")
		agentSecret := r.Header.Get("X-Agent-Secret")
		if agentID == "" || agentSecret == "" {
			http.Error(w, `{"error":"missing agent credentials"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKey("agent_id"), agentID)
		ctx = context.WithValue(ctx, ctxKey("agent_secret"), agentSecret)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
