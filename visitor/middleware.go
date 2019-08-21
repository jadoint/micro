package visitor

import (
	"context"
	"net/http"
	"os"

	"github.com/jadoint/micro/auth"
)

// Middleware adds a Visitor struct to Request
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(os.Getenv("COOKIE_NAME"))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		v := &Visitor{}
		shortToken := cookie.Value
		td, err := auth.ParseToken(shortToken)
		if err == nil {
			v = &Visitor{
				ID:   td.ID,
				Name: td.Name,
			}
		}

		ctx := context.WithValue(r.Context(), GetContextKey(), v)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
