package route

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	"github.com/jadoint/micro/pkg/blog"
	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/visitor"
)

// TagRouter handles all requests to /blog/tag
func TagRouter(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	r.Get("/frequent", func(w http.ResponseWriter, r *http.Request) {
		tags, _ := blog.GetFrequentTags(clients)

		res, err := json.Marshal(struct {
			FrequentTags []*string `json:"frequentTags,omitempty"`
		}{tags})
		logger.HandleError(err)
		w.Write(res)
	})

	r.Get("/{idBlog:[0-9]+}/{jsonName:[a-z0-9_.]+}", func(w http.ResponseWriter, r *http.Request) {
		idBlogParam, err := strconv.Atoi(chi.URLParam(r, "idBlog"))
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		idBlog := int64(idBlogParam)

		tagCsv, _ := blog.GetTagsCSV(clients, idBlog)
		var tags []string
		if tagCsv != "" {
			tags = strings.Split(tagCsv, ",")
		}

		res, err := json.Marshal(struct {
			Tags []string `json:"tags,omitempty"`
		}{tags})
		logger.HandleError(err)
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

		var t blog.Tag
		err = d.Decode(&t)
		logger.HandleError(err)

		// Validation
		err = t.Validate()
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save
		idTag, err := blog.AddTag(clients, idBlog, &t)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := json.Marshal(struct {
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
		err = blog.DeleteTag(clients, idBlog, tag)
		if err != nil {
			errutil.Send(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := json.Marshal(struct {
			ID  int64  `json:"idPost"`
			Tag string `json:"tag"`
		}{idBlog, tag})
		w.Write(res)
	})

	return r
}
