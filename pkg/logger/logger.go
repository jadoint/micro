package logger

import (
	"log"
	"runtime"
)

// Panic logs caller filename and line number and
// sends an Internal Server error code to the user.
func Panic(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if len(v) > 0 {
		log.Panic(file, ":", line, " ", v)
	} else {
		log.Panic(file, ":", line)
	}
}

// HandleError will log error and panic
// if an error is found.
func HandleError(e error) {
	if e == nil {
		return
	}
	Panic(e.Error())
}
