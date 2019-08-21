package auth

import (
	"net/http"
	"time"
)

// AddCookie adds cookie to a visitor's browser
func AddCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(1, 0, 0)
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expire,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

// RemoveCookie removes specific cookie from a visitor's browser
func RemoveCookie(w http.ResponseWriter, name string) {
	expire := time.Now().AddDate(-1, 0, 0)
	cookie := http.Cookie{
		Name:    name,
		Value:   "",
		Path:    "/",
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}
