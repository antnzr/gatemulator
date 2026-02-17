package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/service"
)

type SubscriberKey string

func AuthMiddleware(subscriberService service.SubscriberService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authorization header missing", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")

			subscriber, err := subscriberService.GetOneByToken(token)
			if err != nil || subscriber == nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), SubscriberKey("subscriber"), subscriber)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
