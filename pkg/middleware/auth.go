package middleware

import (
	"context"
	"finance-app/configs"
	"finance-app/pkg/er"
	"finance-app/pkg/jwt"
	"net/http"
	"strings"
)

type UserData struct {
	Id uint
	Email string
	Role string
}

type key string

const (
	ContextUserDataKey key = "ContextUserDataKey"
)

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authedHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authedHeader, "Bearer ") {
			writeUnauthed(w)
			return
		}
		token := strings.TrimPrefix(authedHeader, "Bearer ")
		isValid, data := jwt.NewJwt(config.Auth.Secret).Parse(token)
		if !isValid {
			writeUnauthed(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserDataKey, UserData{
			Id: data.Id,
			Email: data.Email,
			Role: data.Role,
		})
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	http.Error(w, er.ErrNotAuthorized, http.StatusBadRequest)
}
