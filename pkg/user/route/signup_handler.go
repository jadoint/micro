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

func signup(w http.ResponseWriter, r *http.Request, clients *conn.Clients) {
	visitor := visitor.GetVisitor(r)
	if visitor.ID > 0 {
		errutil.Send(w, "Already logged in", http.StatusForbidden)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var ur user.Registration
	err := d.Decode(&ur)
	if err != nil {
		logger.Panic(err.Error())
	}

	err = validate.Struct(ur)
	if err != nil {
		errutil.Send(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if username is unique
	u, _ := user.GetUserByUsername(clients, ur.Username)
	if ur.Username == u.Username {
		errutil.Send(w, "Username already exists", http.StatusForbidden)
		return
	}

	idUser, err := user.AddUser(clients, &ur)
	if err != nil {
		logger.Panic(err.Error())
	}

	// JWT
	tokenString, err := auth.MakeAuthToken(idUser, ur.Username)
	if err != nil {
		logger.Panic(err.Error())
	}

	// Cookie
	auth.AddCookie(w, os.Getenv("COOKIE_SESSION_NAME"), tokenString)

	// Response
	newUser := &user.User{ID: idUser, Username: ur.Username}
	res, err := json.Marshal(newUser)
	if err != nil {
		logger.Panic(err.Error())
	}

	w.Write(res)
}
