package route

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/microcosm-cc/bluemonday"

	"github.com/jadoint/micro/blog/model"
	"github.com/jadoint/micro/conn"
	"github.com/jadoint/micro/errutil"
	"github.com/jadoint/micro/logger"
	"github.com/jadoint/micro/now"
	"github.com/jadoint/micro/validate"
	"github.com/jadoint/micro/visitor"
	"github.com/jadoint/micro/words"
)

// BlogRouter handles all requests to /blog/
func BlogRouter(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	r.Get("/{idBlog}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		idBlog := int64(idBlogParam)

		visitor := visitor.GetVisitor(r)

		blog, err := model.GetBlog(clients, idBlog)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		res, err := json.Marshal(struct {
			*model.Blog
			IDVisitor int64 `json:"idVisitor,omitempty"`
		}{
			blog,
			visitor.ID,
		})
		if err != nil {
			logger.Panic(err.Error(), "Get Blog ID", idBlog)
		}

		w.Write(res)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		visitor := visitor.GetVisitor(r)
		if visitor.ID == 0 {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// Unmarshalling
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()

		var b model.Blog
		err := d.Decode(&b)
		if err != nil {
			logger.Panic(err.Error())
		}
		b.IDAuthor = visitor.ID
		// Sanitize inputs against XSS attacks
		strict := bluemonday.StrictPolicy()
		b.Title = strict.Sanitize(b.Title)
		ugc := bluemonday.UGCPolicy()
		b.Post = ugc.Sanitize(b.Post)
		b.WordCount = words.Count(&b.Post)

		// Validation
		err = validate.Struct(b)
		if err != nil {
			logger.Panic(err.Error())
		}

		// Save
		idBlog, err := model.AddBlog(clients, &b)
		if err != nil {
			logger.Panic(err.Error(), "Add Blog ID", idBlog)
		}

		res, err := json.Marshal(struct {
			ID int64 `json:"idPost"`
		}{idBlog})
		w.Write(res)
	})

	r.Put("/{idBlog}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		idBlog := int64(idBlogParam)

		// Authorization
		v := visitor.GetVisitor(r)
		isAuthorized, status := isAuthorized(clients, v, idBlog)
		if !isAuthorized {
			http.Error(w, "", status)
			return
		}

		// Unmarshalling
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()

		var b model.Blog
		err = d.Decode(&b)
		if err != nil {
			logger.Panic(err.Error())
		}
		b.ID = idBlog
		b.IDAuthor = v.ID
		// Sanitize inputs against XSS attacks
		strict := bluemonday.StrictPolicy()
		b.Title = strict.Sanitize(b.Title)
		ugc := bluemonday.UGCPolicy()
		b.Post = ugc.Sanitize(b.Post)
		b.WordCount = words.Count(&b.Post)
		b.Modified = now.MySQLUTC()

		// Validation
		err = validate.Struct(b)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save
		err = model.UpdateBlog(clients, &b)
		if err != nil {
			logger.Panic(err.Error(), "Update Blog ID", idBlog)
		}

		res, err := json.Marshal(struct {
			ID int64 `json:"idPost"`
		}{idBlog})
		w.Write(res)
	})

	r.Delete("/{idBlog}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		idBlog := int64(idBlogParam)

		// Authorization
		v := visitor.GetVisitor(r)
		isAuthorized, status := isAuthorized(clients, v, idBlog)
		if !isAuthorized {
			http.Error(w, "", status)
			return
		}

		// Delete
		err = model.DeleteBlog(clients, idBlog)
		if err != nil {
			logger.Panic(err.Error(), "Delete Blog ID", idBlog)
		}

		res, err := json.Marshal(struct {
			ID int64 `json:"idPost"`
		}{idBlog})
		w.Write(res)
	})

	// Blog views
	r.Get("/views/{idBlog}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		idBlog := int64(idBlogParam)

		views, err := model.IncrViews(clients, idBlog)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		res, err := json.Marshal(struct {
			Views int64 `json:"views"`
		}{views})
		if err != nil {
			logger.Panic(err.Error(), "Get Blog views", idBlog)
		}

		w.Write(res)
	})

	return r
}

// isAuthorized checks if visitor is authorized to do an action
func isAuthorized(clients *conn.Clients, v *visitor.Visitor, idBlog int64) (bool, int) {
	if v.ID == 0 {
		return false, http.StatusUnauthorized
	}

	idAuthor, err := model.GetIDAuthor(clients, idBlog)
	if err != nil {
		return false, http.StatusNotFound
	}
	if idAuthor != v.ID {
		return false, http.StatusForbidden
	}

	return true, http.StatusOK
}
