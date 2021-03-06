package user

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/cookie"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/hash"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/token"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
)

func login(w http.ResponseWriter, r *http.Request, clients *conn.Clients) {
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
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Validation
	err = validate.Struct(ul)
	if err != nil {
		errutil.Send(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authentication
	u, _ := GetUserByUsername(clients, ul.Username)
	if ul.Username != u.Username {
		errutil.Send(w, "Username and password do not match", http.StatusUnauthorized)
		return
	}

	isMatchingPasswords, err := hash.VerifyPassword(ul.Password, u.Password)
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if !isMatchingPasswords {
		errutil.Send(w, "Username and password do not match", http.StatusUnauthorized)
		return
	}

	// JWT
	dataClaim := visitor.GetVisitorTokenDataClaim(u.ID, u.Username)
	tokenString, err := token.Create(dataClaim)
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Cookie
	cookie.Add(w, os.Getenv("COOKIE_SESSION_NAME"), tokenString)

	// Response
	u.Created = ""
	res, err := json.Marshal(u)
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Write(res)
}
