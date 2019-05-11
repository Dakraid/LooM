package gui

import (
	"fmt"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
	"github.com/hpcloud/tail"

	"github.com/dakraid/LooM/version"
)

var (
	logwin      *ui.Window
	Text        *ui.MultilineEntry
)

func ReadLogs() {
	t, err := tail.TailFile("output.log", tail.Config{Follow: false})

	if err != nil {
		logger.Errorf("Failed to tail logs: %v",err)
	}

	ui.QueueMain(func() {
		for line := range t.Lines {
			Text.Append(line.Text + "\n")
		}
	})
}

func setupLogs() {
	logger.Info("Preparing the logs window")
	logwin = ui.NewWindow(fmt.Sprintf("Loot Master v%s - Logs", version.Version), 480, 300, false)
	logwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		logwin.Destroy()
		return true
	})

	Text = ui.NewMultilineEntry()
	Text.SetReadOnly(true)
	logwin.SetChild(Text)
	logwin.SetMargined(true)

	logwin.Show()
}

func ShowLogs() {
	ui.Main(setupLogs)
	// go ReadLogs()
}