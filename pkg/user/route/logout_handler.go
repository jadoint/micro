package route

import (
	"net/http"
	"os"

	"github.com/jadoint/micro/pkg/cookie"
	"github.com/jadoint/micro/pkg/msg"
)

func logout(w http.ResponseWriter, r *http.Request) {
	cookie.Remove(w, os.Getenv("COOKIE_SESSION_NAME"))

	res := msg.MakeAppMsg("Logged out")
	w.Write(res)
}
