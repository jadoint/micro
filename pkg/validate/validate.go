package validate

import (
	"github.com/go-playground/validator/v10"

	"github.com/jadoint/micro/pkg/errutil"
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
