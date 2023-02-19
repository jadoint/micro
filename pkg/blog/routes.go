package blog

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/jadoint/micro/pkg/clean"
	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/fmtdate"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/paginate"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
	"github.com/jadoint/micro/pkg/words"
)

// RouteBlog handles all requests to /blog/
func RouteBlog(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	r.Get("/{idBlog:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		idBlog := int64(idBlogParam)

		v := visitor.GetVisitor(r)

		b, err := GetPostInit(clients, idBlog)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		// Format date/time
		t, _ := time.Parse("2006-01-02 15:04:05", b.Modified)
		b.Modified = t.Format("January 02, 2006")
		b.ModifiedDatetime = t.Format("20060102150405")

		// Authorization
		if b.IsDraft && b.IDAuthor != v.ID {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		res, err := json.Marshal(struct {
			*Blog
			IDVisitor int64 `json:"idVisitor,omitempty"`
		}{
			b,
			v.ID,
		})
		if err != nil {
			logger.Log(err, "Get Blog ID: "+strconv.FormatInt(idBlog, 10))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Write(res)
	})

	r.Get("/{idBlog:[0-9]+}/{jsonName:[a-z0-9_.]+}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		idBlog := int64(idBlogParam)

		b, err := GetPost(clients, idBlog)
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

		res, err := json.Marshal(struct{ *Blog }{b})
		if err != nil {
			logger.Log(err, "Get Blog ID: "+strconv.FormatInt(idBlog, 10))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Write(res)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		v := visitor.GetVisitor(r)
		if v.ID == 0 {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// Unmarshalling
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()

		var b Blog
		err := d.Decode(&b)
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		b.IDAuthor = v.ID
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
		idBlog, err := Add(clients, &b)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		res, _ := json.Marshal(struct {
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

		var b Blog
		err = d.Decode(&b)
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
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
		b.Modified = fmtdate.MySQLUTCNow()
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
		err = Update(clients, &b)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		res, _ := json.Marshal(struct {
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
		err = Delete(clients, idBlog)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		res, _ := json.Marshal(struct {
			ID int64 `json:"idPost"`
		}{idBlog})
		w.Write(res)
	})

	// Blog views
	r.Put("/views/{idBlog:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		idBlog := int64(idBlogParam)

		views, err := IncrViews(clients, idBlog)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		res, err := json.Marshal(struct {
			Views int64 `json:"views"`
		}{views})
		if err != nil {
			logger.Log(err, "Get Blog views: "+strconv.FormatInt(idBlog, 10))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Write(res)
	})

	r.Get("/latest", func(w http.ResponseWriter, r *http.Request) {
		v := visitor.GetVisitor(r)

		pageNum := paginate.GetPageNum(r)

		var blogs []*Blog
		var err error
		tagParam := r.URL.Query().Get("tag")
		if tagParam != "" {
			// Latest blog listings by tag
			var t Tag
			t.Tag = tagParam

			// Validation
			err = t.Validate()
			if err != nil {
				errutil.Send(w, err.Error(), http.StatusBadRequest)
				return
			}

			blogs, err = GetLatestByTag(clients, t.Tag, pageNum, 10)
			if err != nil {
				http.Error(w, "", http.StatusNotFound)
				return
			}
		} else {
			// Latest blog listings
			blogs, err = GetLatest(clients, pageNum, 10)
		}

		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		res, err := json.Marshal(struct {
			Listings  []*Blog `json:"listings"`
			PageNum   int     `json:"pageNum,omitempty"`
			IDVisitor int64   `json:"idVisitor,omitempty"`
		}{
			Listings:  blogs,
			PageNum:   pageNum,
			IDVisitor: v.ID,
		})
		if err != nil {
			logger.Log(err, "Latest blogs")
			http.Error(w, "", http.StatusBadRequest)
			return
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

		blogs, err := GetRecentAuthorBlogs(clients, idAuthor)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		res, err := json.Marshal(struct {
			Listings  []*Blog `json:"listings"`
			IDVisitor int64   `json:"idVisitor,omitempty"`
		}{
			Listings:  blogs,
			IDVisitor: v.ID,
		})
		if err != nil {
			logger.Log(err, "Recent blogs by author ID: "+strconv.FormatInt(idAuthor, 10))
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		w.Write(res)
	})

	return r
}

// RouteTag handles all requests to /blog/tag
func RouteTag(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	r.Get("/frequent", func(w http.ResponseWriter, r *http.Request) {
		tags, _ := GetFrequentTags(clients)

		res, err := json.Marshal(struct {
			FrequentTags []*string `json:"frequentTags,omitempty"`
		}{tags})
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		w.Write(res)
	})

	r.Get("/{idBlog:[0-9]+}/{jsonName:[a-z0-9_.]+}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		idBlog := int64(idBlogParam)

		tagCsv, _ := GetTagsCSV(clients, idBlog)
		var tags []string
		if tagCsv != "" {
			tags = strings.Split(tagCsv, ",")
		}

		res, err := json.Marshal(struct {
			Tags []string `json:"tags,omitempty"`
		}{tags})
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		w.Write(res)
	})

	r.Post("/{idBlog:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
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

		var t Tag
		err = d.Decode(&t)
		if err != nil {
			logger.Log(err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// Validation
		err = t.Validate()
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save
		idTag, err := AddTag(clients, idBlog, &t)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, _ := json.Marshal(struct {
			ID int64 `json:"idTag"`
		}{idTag})
		w.Write(res)
	})

	r.Delete("/{idBlog:[0-9]+}/{tag}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		idBlog := int64(idBlogParam)
		tag := chi.URLParam(r, "tag")

		// Authorization
		v := visitor.GetVisitor(r)
		isAuthorized, status := isAuthorized(clients, v, idBlog)
		if !isAuthorized {
			http.Error(w, "", status)
			return
		}

		// Delete
		err = DeleteTag(clients, idBlog, tag)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, _ := json.Marshal(struct {
			ID  int64  `json:"idPost"`
			Tag string `json:"tag"`
		}{idBlog, tag})
		w.Write(res)
	})

	return r
}
