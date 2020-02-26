package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

type action int

const (
	logPanic action = iota + 1
	logFatal
	logOnly
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

// writeLogMessage to local file
func writeLogMessage(logType action, v ...interface{}) {
	dir := "logs"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0740)
	}
	t := time.Now()
	filename := t.Format("2006-01-02") + ".log"
	filepath := dir + "/" + filename
	errFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
	// Write error message to local log file
	switch logType {
	case logPanic:
		log.Panic(errMsg)
	case logFatal:
		log.Fatal(errMsg)
	case logOnly:
		log.Println(errMsg)
	default:
		log.Println(errMsg)
	}
}

// logAndPanic logs error and panics.
// Note: It would be a mistake to move this code to Panic()
// and call Panic() in HandleError() because runtime.Caller()
// gets the line number relative to the function that called it.
// Calling Panic() in HandleError() where Panic() has declared
// runtime.Caller(1), would only give HandleError() the line
// number in Panic() and not where the actual error occurred.
func logAndPanic(v ...interface{}) {
	writeLogMessage(logPanic, v)
}

// logAndFatal logs error and calls Fatal.
func logAndFatal(v ...interface{}) {
	writeLogMessage(logFatal, v)
}

// logAndContinue only logs the error.
func logAndContinue(v ...interface{}) {
	writeLogMessage(logOnly, v)
}
