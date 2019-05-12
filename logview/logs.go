package logview

import (
	"fmt"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest" // Required for the UI so it can import the CommonControlsV6 used

	"github.com/dakraid/LooM/version"
)

var (
	logwin *ui.Window
	text   *ui.MultilineEntry
)

// AddEntry() is exported so other parts of the program can add their text to the log view
func AddEntry(input string) {
	ui.QueueMain(func() {
		text.Append(input + "\n")
	})
}

// SetupLogs() is the main function that setups the form and returns the window so it can be used in the main thread
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

	text = ui.NewMultilineEntry()
	text.SetReadOnly(true)
	logwin.SetChild(text)
	logwin.SetMargined(true)

	return logwin
}
