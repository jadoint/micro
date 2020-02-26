package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Panic logs caller filename and line number and
// sends an Internal Server error code to the user.
func Panic(v ...interface{}) {
	logAndPanic(v)
}

// Fatal logs caller filename and line number and
// sends an Internal Server error code to the user.
func Fatal(v ...interface{}) {
	logAndFatal(v)
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

// LogError will only log the error
func LogError(err error) {
	if err == nil {
		return
	}

	errMsg := err.Error()
	logAndContinue(errMsg)
}

// writeLogMessage to local file and return
// the written error message.
func writeLogMessage(v ...interface{}) string {
	t := time.Now()
	filename := t.Format("2006-01-02") + ".log"
	errFile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err.Error())
	}
	defer errFile.Close()
	log.SetOutput(errFile)

	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "go"
	}
	_, file, line, _ := runtime.Caller(3)
	// Timestamp, Docker hostname, error message
	errMsg := t.String() + "," + hostname + "," + file + ": " + strconv.Itoa(line)
	if len(v) > 0 {
		// errMsg + custom error message
		errMsg += " " + fmt.Sprintf("%v", v)
	}
	// Write error to console to help with checking
	// errors in real-time.
	fmt.Println(errMsg)
	return errMsg
}

// logAndPanic logs error and panics.
// Note: It would be a mistake to move this code to Panic()
// and call Panic() in HandleError() because runtime.Caller()
// gets the line number relative to the function that called it.
// Calling Panic() in HandleError() where Panic() has declared
// runtime.Caller(1), would only give HandleError() the line
// number in Panic() and not where the actual error occurred.
func logAndPanic(v ...interface{}) {
	errMsg := writeLogMessage(v)
	log.Panic(errMsg)
}

// logAndFatal logs error and calls Fatal.
func logAndFatal(v ...interface{}) {
	errMsg := writeLogMessage(v)
	log.Fatal(errMsg)
}

// logAndContinue only logs the error.
func logAndContinue(v ...interface{}) {
	errMsg := writeLogMessage(v)
	log.Println(errMsg)
}
