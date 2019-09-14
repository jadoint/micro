package route

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/visitor"
)

// AuthRouter handles signups
func AuthRouter(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	// Rate limiter: first argument is "x requests / second" per IP
	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	r.Use(tollbooth_chi.LimitHandler(lmt))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		v := visitor.GetVisitor(r)

		res, err := json.Marshal(v)
		logger.HandleError(err)

		w.Write(res)
	})

	r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		signup(w, r, clients)
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		login(w, r, clients)
	})

	r.Post("/new-password", func(w http.ResponseWriter, r *http.Request) {
		newPassword(w, r, clients)
	})

	r.Post("/logout", logout)

	return r
}
