package route

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/jadoint/micro/auth"
	"github.com/jadoint/micro/conn"
	"github.com/jadoint/micro/errutil"
	"github.com/jadoint/micro/logger"
	"github.com/jadoint/micro/user/model"
	"github.com/jadoint/micro/validate"
	"github.com/jadoint/micro/visitor"
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

	var ul model.UserLogin
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
	u, _ := model.GetUserByUsername(clients, ul.Username)
	if ul.Username != u.Username {
		errutil.Send(w, "Username does not exist", http.StatusUnauthorized)
		return
	}

	isMatchingPasswords, err := auth.VerifyPasswordHash(ul.Password, u.Password)
	if err != nil {
		logger.Panic(err.Error())
	}
	if !isMatchingPasswords {
		errutil.Send(w, "Password does not match", http.StatusUnauthorized)
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
