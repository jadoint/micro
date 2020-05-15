package validate

import (
	"strconv"
	"strings"

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

// ValidCSVIDs determines if a CSV of integer IDs is valid
func ValidCSVIDs(csv string) bool {
	if csv == "" {
		return false
	}
	ids := strings.Split(csv, ",")
	for _, id := range ids {
		_, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return false
		}
	}
	return true
}
