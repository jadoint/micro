package user

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/jadoint/micro/pkg/cookie"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/hash"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/token"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
)

func (env *Env) login(w http.ResponseWriter, r *http.Request) {
	v := visitor.GetVisitor(r)
	if v.ID > 0 {
		errutil.Send(w, "Already logged in", http.StatusForbidden)
		return
	}

	// Unmarshalling
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var ul Login
	err := d.Decode(&ul)
	logger.HandleError(err)

	// Validation
	err = validate.Struct(ul)
	if err != nil {
		errutil.Send(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authentication
	u, _ := env.GetUserByUsername(ul.Username)
	if ul.Username != u.Username {
		errutil.Send(w, "Username and password do not match", http.StatusUnauthorized)
		return
	}

	isMatchingPasswords, err := hash.VerifyPassword(ul.Password, u.Password)
	logger.HandleError(err)
	if !isMatchingPasswords {
		errutil.Send(w, "Username and password do not match", http.StatusUnauthorized)
		return
	}

	// JWT
	dataClaim := visitor.GetVisitorTokenDataClaim(u.ID, u.Username)
	tokenString, err := token.Create(dataClaim)
	logger.HandleError(err)

	// Cookie
	cookie.Add(w, os.Getenv("COOKIE_SESSION_NAME"), tokenString)

	// Response
	u.Created = ""
	res, err := json.Marshal(u)
	logger.HandleError(err)

	w.Write(res)
}
