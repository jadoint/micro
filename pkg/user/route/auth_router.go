package route

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"

	"github.com/jadoint/micro/pkg/conn"
)

// AuthRouter handles signups
func AuthRouter(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	// Rate limiter: first argument is "x requests / second" per IP
	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	r.Use(tollbooth_chi.LimitHandler(lmt))

	r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		signup(w, r, clients)
	})
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		login(w, r, clients)
	})
	r.Post("/logout", logout)

	return r
}
