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

func signup(w http.ResponseWriter, r *http.Request, clients *conn.Clients) {
	visitor := visitor.GetVisitor(r)
	if visitor.ID > 0 {
		errutil.Send(w, "Already logged in", http.StatusForbidden)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var ur model.UserRegistration
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
	u, _ := model.GetUserByUsername(clients, ur.Username)
	if ur.Username == u.Username {
		errutil.Send(w, "Username already exists", http.StatusForbidden)
		return
	}

	idUser, err := model.AddUser(clients, &ur)
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
	newUser := &model.User{ID: idUser, Username: ur.Username}
	res, err := json.Marshal(newUser)
	if err != nil {
		logger.Panic(err.Error())
	}

	w.Write(res)
}
