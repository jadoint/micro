package logger

import (
	"log"
	"runtime"
)

// Panic logs caller filename and line number and
// sends an Internal Server error code to the user.
func Panic(v ...interface{}) {
	logAndPanic(v)
}

// HandleError will log error and panic
// if an error is found.
func HandleError(err error) {
	if err == nil {
		return
	}

	errMsg := err.Error()
	logAndPanic(errMsg)
}

// logAndPanic logs error and panics.
// Note: It would be a mistake to move this code to Panic()
// and call Panic() in HandleError() because runtime.Caller()
// gets the line number relative to the function that called it.
// Calling Panic() in HandleError() where Panic() has declared
// runtime.Caller(1), would only give HandleError() the line
// number in Panic() and not where the actual error occurred.
func logAndPanic(v ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	if len(v) > 0 {
		log.Panic(file, ":", line, " ", v)
	} else {
		log.Panic(file, ":", line)
	}
}
