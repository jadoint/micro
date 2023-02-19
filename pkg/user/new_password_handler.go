package user

import (
	"encoding/json"
	"net/http"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/hash"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
)

func newPassword(w http.ResponseWriter, r *http.Request, clients *conn.Clients) {
	v := visitor.GetVisitor(r)
	if v.ID == 0 {
		errutil.Send(w, "Not logged in", http.StatusForbidden)
		return
	}

	// Unmarshalling
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var pc PasswordChange
	err := d.Decode(&pc)
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Validation
	err = validate.Struct(pc)
	if err != nil {
		errutil.Send(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authentication
	u, _ := GetUserByUsername(clients, v.Name)
	isMatchingPasswords, err := hash.VerifyPassword(pc.OldPassword, u.Password)
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if !isMatchingPasswords {
		errutil.Send(w, "Password is incorrect", http.StatusForbidden)
		return
	}

	// Save new password
	_ = ChangePassword(clients, v.ID, pc.NewPassword)

	// Response
	res, err := json.Marshal(struct {
		ID int64 `json:"id"`
	}{v.ID})
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Write(res)
}
