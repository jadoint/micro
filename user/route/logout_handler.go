package route

import (
	"net/http"
	"os"

	"github.com/jadoint/micro/auth"
	"github.com/jadoint/micro/msg"
)

func logout(w http.ResponseWriter, r *http.Request) {
	auth.RemoveCookie(w, os.Getenv("COOKIE_SESSION_NAME"))

	res := msg.MakeAppMsg("Logged out")
	w.Write(res)
}
