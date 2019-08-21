package route

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/jadoint/micro/blog/model"
	"github.com/jadoint/micro/conn"
	"github.com/jadoint/micro/logger"
	"github.com/jadoint/micro/visitor"
)

// BlogsRouter handles all requests to /blogs/
func BlogsRouter(clients *conn.Clients) chi.Router {
	r := chi.NewRouter()

	r.Get("/recent/{idAuthor}", func(w http.ResponseWriter, r *http.Request) {
		idAuthorParam, err := strconv.Atoi(chi.URLParam(r, "idAuthor"))
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		idAuthor := int64(idAuthorParam)

		v := visitor.GetVisitor(r)

		blogs, err := model.GetRecentAuthorBlogs(clients, idAuthor)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		res, err := json.Marshal(struct {
			Blogs     []*model.Blog `json:"blogs"`
			IDVisitor int64         `json:"idVisitor,omitempty"`
		}{
			Blogs:     blogs,
			IDVisitor: v.ID,
		})
		if err != nil {
			logger.Panic(err.Error(), "Recent blogs by author ID", idAuthor)
		}

		w.Write(res)
	})

	return r
}
