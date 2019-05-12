package clog

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/logger"

	"github.com/dakraid/LooM/logview"
)

const logPath = "output.log"

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

// InitLogger The cLogger package is just a proxy to the Google/Logger package and has to initialize it first to be able to use it
func InitLogger() {
	flag.Parse()

	err := os.Remove(logPath)
	if err != nil {
		logger.Errorf("Failed to clear log file: %v", err)
	}

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	defer lf.Close()
	logger.Init("OutputLog", *verbose, true, lf)
}

// Info log call
func Info(v ...interface{}) {
	logview.AddEntry(fmt.Sprint(v...))
	logger.Info(v)
}

// Warning log call
func Warning(v ...interface{}) {
	logview.AddEntry(fmt.Sprint(v...))
	logger.Warning(v)
}

// Error log call
func Error(v ...interface{}) {
	logview.AddEntry(fmt.Sprint(v...))
	logger.Error(v)
}

// Fatal log call
func Fatal(v ...interface{}) {
	logview.AddEntry(fmt.Sprint(v...))
	logger.Fatal(v)
}

// Infof Formated Info log call
func Infof(format string, v ...interface{}) {
	logview.AddEntry(fmt.Sprintf(format, v...))
	logger.Infof(format, v)
}

// Warningf Formated Warning log call
func Warningf(format string, v ...interface{}) {
	logview.AddEntry(fmt.Sprintf(format, v...))
	logger.Warningf(format, v)
}

// Errorf Formated Error log call
func Errorf(format string, v ...interface{}) {
	logview.AddEntry(fmt.Sprintf(format, v...))
	logger.Errorf(format, v)
}

// Fatalf Formated Fatal log call
func Fatalf(format string, v ...interface{}) {
	logview.AddEntry(fmt.Sprintf(format, v...))
	logger.Fatalf(format, v)
}
