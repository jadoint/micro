package user

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/jadoint/micro/pkg/conn"
	"github.com/jadoint/micro/pkg/cookie"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/token"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
)

func signup(w http.ResponseWriter, r *http.Request, clients *conn.Clients) {
	v := visitor.GetVisitor(r)
	if v.ID > 0 {
		errutil.Send(w, "Already logged in", http.StatusForbidden)
		return
	}

	isSignupRestricted, err := strconv.ParseBool(os.Getenv("SIGNUPS_RESTRICTED"))
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if isSignupRestricted {
		errutil.Send(w, "Signups no longer accepted", http.StatusForbidden)
		return
	}

	// Marshalling
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	var ur Registration
	err = d.Decode(&ur)
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Validation
	err = validate.Struct(ur)
	if err != nil {
		errutil.Send(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify Recaptcha token
	rr := validateRecaptcha(ur.RecaptchaToken)
	scoreThreshold, err := strconv.ParseFloat(os.Getenv("RECAPTCHA_SCORE_THRESHOLD"), 64)
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if !rr.Success || rr.Score < scoreThreshold {
		errutil.Send(w, "Unable to sign up due to captcha failure. Please refresh the page to try again.", http.StatusForbidden)
		return
	}

	// Check if username is unique
	u, _ := GetUserByUsername(clients, ur.Username)
	if ur.Username == u.Username {
		errutil.Send(w, "Username already exists", http.StatusForbidden)
		return
	}

	// Success: Add user
	idUser, err := AddUser(clients, &ur, rr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// JWT
	dataClaim := visitor.GetVisitorTokenDataClaim(idUser, ur.Username)
	tokenString, err := token.Create(dataClaim)
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Cookie
	cookie.Add(w, os.Getenv("COOKIE_SESSION_NAME"), tokenString)

	// Response
	newUser := &User{ID: idUser, Username: ur.Username}
	res, err := json.Marshal(newUser)
	if err != nil {
		logger.Log(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Write(res)
}

func validateRecaptcha(token string) *RecaptchaResponse {
	captchaFields := url.Values{
		"secret":   {os.Getenv("RECAPTCHA_KEY")},
		"response": {token},
	}
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", captchaFields)
	if err != nil {
		logger.Log(err, "Unable to verify reCaptcha token")
		return nil
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	var rr RecaptchaResponse
	err = d.Decode(&rr)
	if err != nil {
		logger.Log(err)
		return nil
	}

	return &rr
}
