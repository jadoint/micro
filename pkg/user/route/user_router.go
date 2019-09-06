package route

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"

	"github.com/jadoint/micro/pkg/clean"
	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/user"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
)

// UserRouter handles user requests
func UserRouter(clients *conn.Clients) chi.Router {
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

		username, err := user.GetUsername(clients, idUser)
		if err != nil {
			errutil.Send(w, "", http.StatusNotFound)
			return
		}

		// Response
		res, err := json.Marshal(struct {
			Username string `json:"username"`
		}{username})
		if err != nil {
			logger.Panic(err.Error())
		}

		w.Write(res)
	})

	r.Post("/names", func(w http.ResponseWriter, r *http.Request) {
		// Unmarshalling
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()

		var uids user.UserIDs
		err := d.Decode(&uids)
		if err != nil {
			logger.Panic(err.Error())
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

		names, err := user.GetUsernames(clients, &uids)
		if err != nil {
			errutil.Send(w, "", http.StatusNotFound)
			return
		}

		// Response
		res, err := json.Marshal(struct {
			Usernames []*user.Username `json:"usernames"`
		}{names})
		if err != nil {
			logger.Panic(err.Error())
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

		u, err := user.GetAbout(clients, idUser)
		if err != nil {
			errutil.Send(w, "", http.StatusNotFound)
			return
		}

		// Response
		res, err := json.Marshal(struct {
			Title string `json:"title"`
			About string `json:"about"`
		}{
			Title: u.Title,
			About: u.About,
		})
		if err != nil {
			logger.Panic(err.Error())
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

		var u user.User
		err = d.Decode(&u)
		if err != nil {
			logger.Panic(err.Error())
		}
		u.ID = v.ID
		// Strip inputs of all tags
		strict := clean.Strict()
		u.Title = strict.Sanitize(u.Title)
		u.About = strict.Sanitize(u.About)

		// Validation
		err = validate.Struct(u)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save
		err = user.UpdateAbout(clients, &u)
		if err != nil {
			logger.Panic(err.Error(), "Update About", u.ID)
		}

		// Response
		res, err := json.Marshal(struct {
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
		err = user.DeleteAbout(clients, v.ID)
		if err != nil {
			logger.Panic(err.Error(), "Delete About", v.ID)
		}

		// Response
		res, err := json.Marshal(struct {
			ID int64 `json:"idUser"`
		}{idUser})
		w.Write(res)
	})

	return r
}
