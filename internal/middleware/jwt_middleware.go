package middleware

import (
	"context"
	"net/http"
	"petstore/internal/config"

	"github.com/go-chi/jwtauth"
)

func JWTAuthMiddleware(next http.Handler) http.Handler{
	return jwtauth.Verifier(config.TokenAuth)(jwtauth.Authenticator(next))
}

func CetUserFromContext(ctx context.Context) string{
	_, claims, _ := jwtauth.FromContext(ctx)
	if username, ok := claims["username"].(string); ok{
		return username
	}
	return ""
}
