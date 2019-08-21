package route

import (
	"net/http"
	"os"

	"github.com/jadoint/micro/auth"
	"github.com/jadoint/micro/errutil"
	"github.com/jadoint/micro/msg"
	"github.com/jadoint/micro/visitor"
)

func logout(w http.ResponseWriter, r *http.Request) {
	visitor := visitor.GetVisitor(r)
	if visitor.ID <= 0 {
		errutil.Send(w, "Not logged in", http.StatusForbidden)
		return
	}

	auth.RemoveCookie(w, os.Getenv("COOKIE_SESSION_NAME"))

	res := msg.MakeAppMsg("Logged out")
	w.Write(res)
}
