package visitor

import (
	"net/http"
)

// Visitor contains visitor details
type Visitor struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetVisitor returns visitor details from cookie
func GetVisitor(r *http.Request) *Visitor {
	visitor := &Visitor{}
	cookieVisitor := r.Context().Value(GetContextKey())
	if cookieVisitor != nil {
		visitor = cookieVisitor.(*Visitor)
	}
	return visitor
}
