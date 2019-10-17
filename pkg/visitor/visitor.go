package visitor

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/jadoint/micro/pkg/contextkey"
	"github.com/jadoint/micro/pkg/token"
)

// Visitor contains visitor details
type Visitor struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// GetVisitor returns visitor details from cookie
func GetVisitor(r *http.Request) *Visitor {
	visitor := &Visitor{}
	cookieVisitor := r.Context().Value(contextkey.GetVisitorKey())
	if cookieVisitor != nil {
		visitor = cookieVisitor.(*Visitor)
	}
	return visitor
}

// GetVisitorFromCookie returns a Visitor struct
// after parsing the session cookie.
func GetVisitorFromCookie(cookie *http.Cookie) *Visitor {
	v := &Visitor{}
	shortToken := cookie.Value
	claims, err := token.Parse(shortToken)
	// iat := int64(claims["iat"].(float64))
	dataClaim := claims["data"].(map[string]interface{})
	id := int64(dataClaim["id"].(float64))
	name := dataClaim["name"].(string)
	if err == nil {
		v = &Visitor{
			ID:   id,
			Name: name,
		}
	}
	return v
}

// GetVisitorTokenDataClaim returns a jwt.MapClaims type
// for the "data" section of a token.
func GetVisitorTokenDataClaim(id int64, name string) *jwt.MapClaims {
	return &jwt.MapClaims{
		"id":   id,
		"name": name,
	}
}
