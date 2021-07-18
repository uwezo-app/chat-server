package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	c "github.com/uwezo-app/chat-server/controller"
)

type Key string

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var access_token = strings.Split(r.Header.Get("Authorization"), " ")[1]
		access_token = strings.TrimSpace(access_token)

		if access_token == "" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(struct {
				Code    int
				Message string
			}{
				Code:    http.StatusForbidden,
				Message: "Missing Auth Token",
			})
			return
		}

		tk := &c.Token{}

		_, err := jwt.ParseWithClaims(access_token, tk, func(token *jwt.Token) (interface{}, error) {
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
