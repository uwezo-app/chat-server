package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/uwezo-app/chat-server/db"
)

type Key string

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accessToken = strings.Split(r.Header.Get("Authorization"), " ")[1]
		accessToken = strings.TrimSpace(accessToken)

		if accessToken == "" {
			w.WriteHeader(http.StatusForbidden)
			err := json.NewEncoder(w).Encode(struct {
				Code    int
				Message string
			}{
				Code:    http.StatusForbidden,
				Message: "Missing Auth Token",
			})
			if err != nil {
				return
			}
			return
		}

		tk := &db.Token{}

		_, err := jwt.ParseWithClaims(accessToken, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			json.NewEncoder(w).Encode(struct {
				Code    int
				Message string
			}{
				Code:    http.StatusForbidden,
				Message: "Missing Auth Token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), Key("user"), tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
