package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/jadoint/micro/pkg/contextkey"
	"github.com/jadoint/micro/pkg/visitor"
)

// Middleware adds a Visitor struct to Request
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(os.Getenv("COOKIE_SESSION_NAME"))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		v := visitor.GetVisitorFromCookie(cookie)
		ctx := context.WithValue(r.Context(), contextkey.GetVisitorKey(), v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
