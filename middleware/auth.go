package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
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
				Message: "Missing Auth CustomClaims",
			})
			if err != nil {
				return
			}
			return
		}

		tk := &db.CustomClaims{}

		_, err := jwt.ParseWithClaims(accessToken, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			json.NewEncoder(w).Encode(struct {
				Code    int
				Message string
			}{
				Code:    http.StatusForbidden,
				Message: "Missing Auth CustomClaims",
			})
			return
		}

		ctx := context.WithValue(r.Context(), Key("user"), tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
