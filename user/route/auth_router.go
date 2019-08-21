package route

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/jadoint/micro/conn"
)

// AuthRouter handles signups
func AuthRouter(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		signup(w, r, clients)
	})
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		login(w, r, clients)
	})
	r.Post("/logout", logout)

	return r
}
