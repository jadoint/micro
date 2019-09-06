package route

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/jadoint/micro/pkg/auth"
	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/user"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
)

func login(w http.ResponseWriter, r *http.Request, clients *conn.Clients) {
	visitor := visitor.GetVisitor(r)
	if visitor.ID > 0 {
		errutil.Send(w, "Already logged in", http.StatusForbidden)
		return
	}

	// Unmarshalling
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var ul user.Login
	err := d.Decode(&ul)
	if err != nil {
		logger.Panic(err.Error())
	}

	// Validation
	err = validate.Struct(ul)
	if err != nil {
		errutil.Send(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authentication
	u, _ := user.GetUserByUsername(clients, ul.Username)
	if ul.Username != u.Username {
		errutil.Send(w, "Username and password do not match", http.StatusUnauthorized)
		return
	}

	isMatchingPasswords, err := auth.VerifyPasswordHash(ul.Password, u.Password)
	if err != nil {
		logger.Panic(err.Error())
	}
	if !isMatchingPasswords {
		errutil.Send(w, "Username and password do not match", http.StatusUnauthorized)
		return
	}

	// JWT
	tokenString, err := auth.MakeAuthToken(u.ID, u.Username)
	if err != nil {
		logger.Panic(err.Error())
	}

	// Cookie
	auth.AddCookie(w, os.Getenv("COOKIE_SESSION_NAME"), tokenString)

	// Response
	u.Created = ""
	res, err := json.Marshal(u)
	if err != nil {
		logger.Panic(err.Error())
	}

	w.Write(res)
}
