package logview

import (
	"fmt"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	_ "github.com/go-sql-driver/mysql"

	"github.com/dakraid/LooM/version"
)

var (
	logwin      *ui.Window
	Text        *ui.MultilineEntry
)

func AddEntry(input string) {
	ui.QueueMain(func() {
		Text.Append(input + "\n")
	})
}

func SetupLogs() *ui.Window {
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

	return logwin
}