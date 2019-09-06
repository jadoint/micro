package route

import (
	"net/http"
	"os"

	"github.com/jadoint/micro/pkg/auth"
	"github.com/jadoint/micro/pkg/msg"
)

func logout(w http.ResponseWriter, r *http.Request) {
	auth.RemoveCookie(w, os.Getenv("COOKIE_SESSION_NAME"))

	res := msg.MakeAppMsg("Logged out")
	w.Write(res)
}
