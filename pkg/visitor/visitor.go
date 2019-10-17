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
	if err != nil {
		return v
	}

	// iat := int64(claims["iat"].(float64))
	dataClaim := claims["data"].(map[string]interface{})

	// ID
	if idClaim, ok := dataClaim["id"]; ok {
		v.ID = int64(idClaim.(float64))
	}

	// Name
	if name, ok := dataClaim["name"]; ok {
		v.Name = name.(string)
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
