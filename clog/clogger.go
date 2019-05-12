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

func InitLogger() {
	flag.Parse()

	err := os.Remove(logPath)
	if err != nil {
		logger.Fatalf("Failed to clear log file: %v", err)
	}

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	defer lf.Close()
	logger.Init("OutputLog", *verbose, true, lf)
}

func Info(v ...interface{}) {
	logview.AddEntry(fmt.Sprint(v...))
	logger.Info(v)
}

func Warning(v ...interface{}) {
	logview.AddEntry(fmt.Sprint(v...))
	logger.Warning(v)
}

func Error(v ...interface{}) {
	logview.AddEntry(fmt.Sprint(v...))
	logger.Error(v)
}

func Fatal(v ...interface{}) {
	logview.AddEntry(fmt.Sprint(v...))
	logger.Fatal(v)
}

func Infof(format string, v ...interface{}) {
	logview.AddEntry(fmt.Sprintf(format, v...))
	logger.Infof(format,v)
}

func Warningf(format string, v ...interface{}) {
	logview.AddEntry(fmt.Sprintf(format, v...))
	logger.Warningf(format,v)
}

func Errorf(format string, v ...interface{}) {
	logview.AddEntry(fmt.Sprintf(format, v...))
	logger.Errorf(format,v)
}

func Fatalf(format string, v ...interface{}) {
	logview.AddEntry(fmt.Sprintf(format, v...))
	logger.Fatalf(format,v)
}