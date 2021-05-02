package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

type action int

const (
	logPanic action = iota + 1
	logFatal
	logOnly
)

// Log logs error to Stdout
func Log(err error, msg ...string) {
	if err == nil {
		return
	}

	// Docker hostname, caller location, error message
	_, file, line, _ := runtime.Caller(1)
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "go"
	}
	errMsg := hostname + "," + err.Error() + "," + file + ": " + strconv.Itoa(line)
	if len(msg) > 0 {
		var sb strings.Builder
		var msgLen int
		for _, v := range msg {
			msgLen, _ = sb.WriteString(" | " + v)
		}
		if msgLen > 0 {
			errMsg = sb.String()
		}
	}

	// Save log to a file if enabled through
	// a log directory environment variable.
	logsDir := os.Getenv("LOGS_DIR")
	var logFile *os.File
	if logsDir != "" {
		if _, err := os.Stat(logsDir); os.IsNotExist(err) {
			os.Mkdir(logsDir, os.ModePerm)
		}
		t := time.Now()
		filename := t.Format("2006-01-02") + ".log"
		filepath := logsDir + "/" + filename
		logFile, err = os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Print(err.Error())
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}
	if logFile == nil {
		logFile = os.Stdout
	}

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, logFile, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	logger.Error(
		errMsg,
		zap.Error(err))
}

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
		os.Mkdir(dir, os.ModePerm)
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
