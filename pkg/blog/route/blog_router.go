package route

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"github.com/jadoint/micro/pkg/blog"
	"github.com/jadoint/micro/pkg/clean"
	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/now"
	"github.com/jadoint/micro/pkg/paginate"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
	"github.com/jadoint/micro/pkg/words"
)

// BlogRouter handles all requests to /blog/
func BlogRouter(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	r.Get("/{idBlog:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		idBlog := int64(idBlogParam)

		visitor := visitor.GetVisitor(r)

		b, err := blog.Get(clients, idBlog)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		// Format date/time
		t, _ := time.Parse("2006-01-02 15:04:05", b.Created)
		b.Created = t.Format("January 02, 2006")
		t, _ = time.Parse("2006-01-02 15:04:05", b.Modified)
		b.Modified = t.Format("January 02, 2006")
		b.ModifiedDatetime = t.Format("20060102150405")

		// Authorization
		if b.IsDraft && b.IDAuthor != visitor.ID {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		res, err := json.Marshal(struct {
			*blog.Blog
			IDVisitor int64 `json:"idVisitor,omitempty"`
		}{
			b,
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

		var b blog.Blog
		err := d.Decode(&b)
		if err != nil {
			logger.Panic(err.Error())
		}
		b.IDAuthor = visitor.ID
		// Strip title of all tags
		strict := clean.Strict()
		b.Title = strict.Sanitize(b.Title)
		// Sanitize post against XSS attacks
		ugc := clean.UGC()
		b.Post = ugc.Sanitize(b.Post)
		b.WordCount = words.Count(&b.Post)
		// Privacy
		if b.IsDraft {
			b.IsUnlisted = true
		}

		// Validation
		err = validate.Struct(b)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save
		idBlog, err := blog.Add(clients, &b)
		if err != nil {
			logger.Panic(err.Error(), "Add Blog ID", idBlog)
		}

		res, err := json.Marshal(struct {
			ID int64 `json:"idPost"`
		}{idBlog})
		w.Write(res)
	})

	r.Put("/{idBlog:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
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

		var b blog.Blog
		err = d.Decode(&b)
		if err != nil {
			logger.Panic(err.Error())
		}
		b.ID = idBlog
		b.IDAuthor = v.ID
		// Strip title of all tags
		strict := clean.Strict()
		b.Title = strict.Sanitize(b.Title)
		// Sanitize post against XSS attacks
		ugc := clean.UGC()
		b.Post = ugc.Sanitize(b.Post)
		b.WordCount = words.Count(&b.Post)
		b.Modified = now.MySQLUTC()
		// Privacy
		if b.IsDraft {
			b.IsUnlisted = true
		}

		// Validation
		err = validate.Struct(b)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save
		err = blog.Update(clients, &b)
		if err != nil {
			logger.Panic(err.Error(), "Update Blog ID", idBlog)
		}

		res, err := json.Marshal(struct {
			ID int64 `json:"idPost"`
		}{idBlog})
		w.Write(res)
	})

	r.Delete("/{idBlog:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
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
		err = blog.Delete(clients, idBlog)
		if err != nil {
			logger.Panic(err.Error(), "Delete Blog ID", idBlog)
		}

		res, err := json.Marshal(struct {
			ID int64 `json:"idPost"`
		}{idBlog})
		w.Write(res)
	})

	// Blog views
	r.Get("/views/{idBlog:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		idBlog := int64(idBlogParam)

		views, err := blog.IncrViews(clients, idBlog)
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

	r.Get("/latest", func(w http.ResponseWriter, r *http.Request) {
		v := visitor.GetVisitor(r)

		pageNum, err := paginate.GetPageNum(r)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		var blogs []*blog.Blog
		tagParam := r.URL.Query().Get("tag")
		if tagParam != "" {
			// Latest blog listings by tag
			var t blog.Tag
			t.Tag = tagParam

			// Validation
			err = t.Validate()
			if err != nil {
				errutil.Send(w, err.Error(), http.StatusBadRequest)
				return
			}

			blogs, err = blog.GetLatestByTag(clients, t.Tag, pageNum, 10)
			if err != nil {
				http.Error(w, "", http.StatusNotFound)
				return
			}
		} else {
			// Latest blog listings
			blogs, err = blog.GetLatest(clients, pageNum, 10)
		}

		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		res, err := json.Marshal(struct {
			Listings  []*blog.Blog `json:"listings"`
			PageNum   int          `json:"pageNum,omitempty"`
			IDVisitor int64        `json:"idVisitor,omitempty"`
		}{
			Listings:  blogs,
			PageNum:   pageNum,
			IDVisitor: v.ID,
		})
		if err != nil {
			logger.Panic(err.Error(), "Latest blogs")
		}

		w.Write(res)
	})

	r.Get("/recent/{idAuthor:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		idAuthorParam, err := strconv.Atoi(chi.URLParam(r, "idAuthor"))
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		idAuthor := int64(idAuthorParam)

		v := visitor.GetVisitor(r)

		blogs, err := blog.GetRecentAuthorBlogs(clients, idAuthor)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		res, err := json.Marshal(struct {
			Listings  []*blog.Blog `json:"listings"`
			IDVisitor int64        `json:"idVisitor,omitempty"`
		}{
			Listings:  blogs,
			IDVisitor: v.ID,
		})
		if err != nil {
			logger.Panic(err.Error(), "Recent blogs by author ID", idAuthor)
		}

		w.Write(res)
	})

	return r
}
