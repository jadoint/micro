package route

import (
	"net/http"

	"github.com/jadoint/micro/blog"
	"github.com/jadoint/micro/conn"
	"github.com/jadoint/micro/visitor"
)

// isAuthorized checks if visitor is authorized to do an action
func isAuthorized(clients *conn.Clients, v *visitor.Visitor, idBlog int64) (bool, int) {
	if v.ID == 0 {
		return false, http.StatusUnauthorized
	}

	idAuthor, err := blog.GetIDAuthor(clients, idBlog)
	if err != nil {
		return false, http.StatusNotFound
	}
	if idAuthor != v.ID {
		return false, http.StatusForbidden
	}

	return true, http.StatusOK
}
