package errutil

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Errors compiles all errors into a single error
func Errors(v *validator.Validate, err error) error {
	translator := en.New()
	uni := ut.New(translator, translator)

	// This is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		return err
	}

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return err
	}

	var errorMessages []string
	for _, e := range err.(validator.ValidationErrors) {
		errorMessages = append(errorMessages, e.Translate(trans))
	}

	return errors.New(strings.Join(errorMessages, ","))
}

// MakeJSONError converts error message to json format
func MakeJSONError(e error) error {
	errorMessage := struct {
		Error string `json:"error"`
	}{
		Error: e.Error(),
	}

	errorMsgJSON, err := json.Marshal(errorMessage)
	if err != nil {
		return err
	}
	return errors.New(string(errorMsgJSON))
}

// Send sends HTTP error to browser
func Send(w http.ResponseWriter, errMsg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonErr := MakeJSONError(errors.New(errMsg))
	http.Error(w, jsonErr.Error(), status)
}
