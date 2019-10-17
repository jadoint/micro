package route

import (
	"encoding/json"
	"net/http"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/hash"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/user"
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

	var pc user.PasswordChange
	err := d.Decode(&pc)
	logger.HandleError(err)

	// Validation
	err = validate.Struct(pc)
	if err != nil {
		errutil.Send(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authentication
	u, _ := user.GetUserByUsername(clients, v.Name)
	isMatchingPasswords, err := hash.VerifyPassword(pc.OldPassword, u.Password)
	logger.HandleError(err)
	if !isMatchingPasswords {
		errutil.Send(w, "Password is incorrect", http.StatusForbidden)
		return
	}

	// Save new password
	err = user.ChangePassword(clients, v.ID, pc.NewPassword)

	// Response
	res, err := json.Marshal(struct {
		ID int64 `json:"id"`
	}{v.ID})
	logger.HandleError(err)

	w.Write(res)
}
