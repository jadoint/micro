package blog

import (
	"net/http"

	"github.com/jadoint/micro/pkg/visitor"
)

// isAuthorized checks if visitor is authorized to do an action
func (env *Env) isAuthorized(v *visitor.Visitor, idBlog int64) (bool, int) {
	if v.ID == 0 {
		return false, http.StatusUnauthorized
	}

	idAuthor, err := env.GetIDAuthor(idBlog)
	if err != nil {
		return false, http.StatusNotFound
	}
	if idAuthor != v.ID {
		return false, http.StatusForbidden
	}

	return true, http.StatusOK
}
