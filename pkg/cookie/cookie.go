package cookie

import (
	"net/http"
	"os"
	"time"
)

// Add adds cookie to a visitor's browser
func Add(w http.ResponseWriter, name string, value string) {
	isSecureCookie := os.Getenv("IS_SECURE_COOKIE")
	isSecure := false
	if isSecureCookie == "true" {
		isSecure = true
	}
	expire := time.Now().AddDate(1, 0, 0)
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expire,
		Secure:   isSecure,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

// Remove removes specific cookie from a visitor's browser
func Remove(w http.ResponseWriter, name string) {
	expire := time.Now().AddDate(-1, 0, 0)
	cookie := http.Cookie{
		Name:    name,
		Value:   "",
		Path:    "/",
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

// HasCookie determines if a visitor has a specific cookie
func HasCookie(r *http.Request, name string) bool {
	_, err := r.Cookie(name)
	if err != nil {
		return false
	}
	return true
}
