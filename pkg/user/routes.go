package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi/v5"

	"github.com/jadoint/micro/pkg/clean"
	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
)

// RouteAuth handles signups
func RouteAuth(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	// Rate limiter: first argument is "x requests / second" per IP
	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	r.Use(tollbooth_chi.LimitHandler(lmt))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		v := visitor.GetVisitor(r)

		res, err := json.Marshal(v)
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

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

// RouteUser handles user requests
func RouteUser(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	// Rate limiter: first argument is "x requests / second" per IP
	lmt := tollbooth.NewLimiter(100, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"X-Forwarded-For", "RemoteAddr", "X-Real-IP"})
	r.Use(tollbooth_chi.LimitHandler(lmt))

	r.Get("/name/{idUser:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		idUserParam, err := strconv.Atoi(chi.URLParam(r, "idUser"))
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		idUser := int64(idUserParam)

		username, err := GetUsername(clients, idUser)
		if err != nil {
			errutil.Send(w, "", http.StatusNotFound)
			return
		}

		// Response
		res, err := json.Marshal(struct {
			Username string `json:"username"`
		}{username})
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Write(res)
	})

	r.Post("/names", func(w http.ResponseWriter, r *http.Request) {
		// Unmarshalling
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()

		var uids IDs
		err := d.Decode(&uids)
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if len(uids.IDs) == 0 {
			errutil.Send(w, "", http.StatusBadRequest)
			return
		}

		// Validation
		err = validate.Struct(uids)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		names, err := GetUsernames(clients, &uids)
		if err != nil {
			errutil.Send(w, "", http.StatusNotFound)
			return
		}

		// Response
		res, err := json.Marshal(struct {
			Usernames []*Username `json:"usernames"`
		}{names})
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Write(res)
	})

	r.Get("/about/{idUser:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		idUserParam, err := strconv.Atoi(chi.URLParam(r, "idUser"))
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		idUser := int64(idUserParam)

		a, err := GetAbout(clients, idUser)
		if err != nil {
			errutil.Send(w, "", http.StatusNotFound)
			return
		}

		// Response
		res, err := json.Marshal(struct {
			Title string `json:"title"`
			About string `json:"about"`
		}{
			Title: a.Title,
			About: a.About,
		})
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Write(res)
	})

	r.Put("/about/{idUser:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		idUserParam, err := strconv.Atoi(chi.URLParam(r, "idUser"))
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		idUser := int64(idUserParam)

		// Authorization
		v := visitor.GetVisitor(r)
		if v.ID == 0 || v.ID != idUser {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// Unmarshalling
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()

		var a About
		err = d.Decode(&a)
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		// Strip inputs of all tags
		strict := clean.Strict()
		a.Title = strict.Sanitize(a.Title)
		a.About = strict.Sanitize(a.About)

		// Validation
		err = validate.Struct(a)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save
		err = UpdateAbout(clients, v.ID, &a)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// Response
		res, _ := json.Marshal(struct {
			ID int64 `json:"idUser"`
		}{idUser})
		w.Write(res)
	})

	r.Delete("/about/{idUser:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		idUserParam, err := strconv.Atoi(chi.URLParam(r, "idUser"))
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		idUser := int64(idUserParam)

		// Authorization
		v := visitor.GetVisitor(r)
		if v.ID == 0 || v.ID != idUser {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// Delete
		err = DeleteAbout(clients, v.ID)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// Response
		res, _ := json.Marshal(struct {
			ID int64 `json:"idUser"`
		}{idUser})
		w.Write(res)
	})

	return r
}
