package validate

import (
	"gopkg.in/go-playground/validator.v9"

	"github.com/jadoint/micro/errutil"
)

// Struct checks a struct against its validation rules
// and returns a formatted error type if it fails.
func Struct(st interface{}) error {
	v := validator.New()

	err := v.Struct(st)
	if err != nil {
		return errutil.Errors(v, err)
	}
	return nil
}
