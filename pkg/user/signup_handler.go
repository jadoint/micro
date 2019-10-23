package user

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/jadoint/micro/pkg/cookie"
	"github.com/jadoint/micro/pkg/errutil"
	"github.com/jadoint/micro/pkg/logger"
	"github.com/jadoint/micro/pkg/token"
	"github.com/jadoint/micro/pkg/validate"
	"github.com/jadoint/micro/pkg/visitor"
)

func (env *Env) signup(w http.ResponseWriter, r *http.Request) {
	v := visitor.GetVisitor(r)
	if v.ID > 0 {
		errutil.Send(w, "Already logged in", http.StatusForbidden)
		return
	}

	isSignupRestricted, err := strconv.ParseBool(os.Getenv("SIGNUPS_RESTRICTED"))
	logger.HandleError(err)
	if isSignupRestricted {
		errutil.Send(w, "Signups no longer accepted", http.StatusForbidden)
		return
	}

	// Marshalling
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	var ur Registration
	err = d.Decode(&ur)
	logger.HandleError(err)

	// Validation
	err = validate.Struct(ur)
	if err != nil {
		errutil.Send(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verify Recaptcha token
	rr := validateRecaptcha(ur.RecaptchaToken)
	scoreThreshold, err := strconv.ParseFloat(os.Getenv("RECAPTCHA_SCORE_THRESHOLD"), 64)
	logger.HandleError(err)
	if !rr.Success || rr.Score < scoreThreshold {
		errutil.Send(w, "Unable to sign up due to captcha failure. Please refresh the page to try again.", http.StatusForbidden)
		return
	}

	// Check if username is unique
	u, _ := env.GetUserByUsername(ur.Username)
	if ur.Username == u.Username {
		errutil.Send(w, "Username already exists", http.StatusForbidden)
		return
	}

	// Success: Add user
	idUser, err := env.AddUser(&ur, rr)
	logger.HandleError(err)

	// JWT
	dataClaim := visitor.GetVisitorTokenDataClaim(idUser, ur.Username)
	tokenString, err := token.Create(dataClaim)
	logger.HandleError(err)

	// Cookie
	cookie.Add(w, os.Getenv("COOKIE_SESSION_NAME"), tokenString)

	// Response
	newUser := &User{ID: idUser, Username: ur.Username}
	res, err := json.Marshal(newUser)
	logger.HandleError(err)

	w.Write(res)
}

func validateRecaptcha(token string) *RecaptchaResponse {
	captchaFields := url.Values{
		"secret":   {os.Getenv("RECAPTCHA_KEY")},
		"response": {token},
	}
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", captchaFields)
	if err != nil {
		logger.Panic("Unable to verify reCaptcha token: %s", err.Error())
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	var rr RecaptchaResponse
	err = d.Decode(&rr)
	logger.HandleError(err)

	return &rr
}
