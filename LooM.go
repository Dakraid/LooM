package main

import (
	"flag"
	"os"

	"github.com/google/logger"

	"github.com/dakraid/LooM/gui"
	"github.com/dakraid/LooM/version"
)

const logPath = "output.log"

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")

func main() {
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
	defer logger.Init("OutputLog", *verbose, true, lf).Close()

	logger.Infof("Starting Loot Master  v%s", version.Version)

	// TODO: Fix the log window. Right now it won't even show up and with other implementations it doesn't update properly.
	go gui.ShowLogs()
	gui.ShowLogin()
}
