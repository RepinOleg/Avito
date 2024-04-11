package handler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RepinOleg/Banner_service/internal/response"
)

type key string

const userIDKey key = "UserID"

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func (h *Handler) TokenValidationMiddleware(next http.Handler, handlerName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("token")
		if header == "" {
			http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
			return
		}

		var adminFlag bool
		if handlerName != "GetUserBanner" {
			adminFlag = true
		}

		userID, err := h.services.Authorization.ParseToken(headerParts[1], adminFlag)
		if err != nil {
			response.HandleError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
