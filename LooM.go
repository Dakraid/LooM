package main

import (
	"github.com/andlabs/ui"

	"github.com/dakraid/LooM/clog"
	"github.com/dakraid/LooM/gui"
	"github.com/dakraid/LooM/logview"
	"github.com/dakraid/LooM/version"
)

func showWindow() {
	err := ui.Main(func() {
		logview.SetupLogs().Show()
		gui.SetupLogin().Show()
	})

	if err != nil {
		panic(err)
	}
}

func main() {

	clog.InitLogger()
	clog.Infof("Starting Loot Master  v%s", version.Version)

	// TODO: Fix the log window. Right now it won't even show up and with other implementations it doesn't update properly.
	showWindow()
}
