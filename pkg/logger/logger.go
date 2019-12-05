package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
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

// logAndPanic logs error and panics.
// Note: It would be a mistake to move this code to Panic()
// and call Panic() in HandleError() because runtime.Caller()
// gets the line number relative to the function that called it.
// Calling Panic() in HandleError() where Panic() has declared
// runtime.Caller(1), would only give HandleError() the line
// number in Panic() and not where the actual error occurred.
func logAndPanic(v ...interface{}) {
	errFile, err := os.OpenFile("errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err.Error())
	}
	defer errFile.Close()
	log.SetOutput(errFile)

	_, file, line, _ := runtime.Caller(2)
	if len(v) > 0 {
		fmt.Println(file, ":", line, " ", v)
		log.Panic(file, ":", line, " ", v)
	} else {
		fmt.Println(file, ":", line)
		log.Panic(file, ":", line)
	}
}

// logAndFatal logs error and calls Fatal.
func logAndFatal(v ...interface{}) {
	errFile, err := os.OpenFile("errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err.Error())
	}
	defer errFile.Close()
	log.SetOutput(errFile)

	_, file, line, _ := runtime.Caller(2)
	if len(v) > 0 {
		fmt.Println(file, ":", line, " ", v)
		log.Fatal(file, ":", line, " ", v)
	} else {
		fmt.Println(file, ":", line)
		log.Fatal(file, ":", line)
	}
}

// logAndContinue only logs the error.
func logAndContinue(v ...interface{}) {
	errFile, err := os.OpenFile("errors.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err.Error())
	}
	defer errFile.Close()
	log.SetOutput(errFile)

	_, file, line, _ := runtime.Caller(2)
	if len(v) > 0 {
		fmt.Println(file, ":", line, " ", v)
		log.Println(file, ":", line, " ", v)
	} else {
		fmt.Println(file, ":", line)
		log.Println(file, ":", line)
	}
}
