package gui

import (
	"database/sql"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"

	"golang.org/x/crypto/bcrypt"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest" // Required for the UI so it can import the CommonControlsV6 used
	_ "github.com/go-sql-driver/mysql"    // The MySQL driver registers itself as available to the database/sql package

	"github.com/dakraid/LooM/clog"
	"github.com/dakraid/LooM/database"
	"github.com/dakraid/LooM/version"
)

var (
	loginwin   *ui.Window
	dataSource string
)

func openBrowser(url string) {
	clog.Info("Attempting to open link in browser")
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Unsupported Platform %s", runtime.GOOS)
	}
	if err != nil {
		clog.Fatalf("Error while opening URL: %v", err)
	}
}

func hashAndSalt(pwd []byte) string {
	clog.Info("Generating hash and salt for the password")
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)

	if err != nil {
		clog.Fatalf("Error while generating hash: %v", err)
	}

	return string(hash)
}

// TODO: Handle the database connection in its own function or own package
func connectDatabase() *sql.DB {
	clog.Info("Trying to establish database connection")

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		clog.Fatalf("Could not establish database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		clog.Fatalf("Database connectivity issue: %v", err)
	}
	defer db.Close()

	return db
}

func registerAccount(name, hash string) {
	var err error
	clog.Info("Trying to establish database connection")

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		clog.Fatalf("Could not establish database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		clog.Fatalf("Database connectivity issue: %v", err)
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO Accounts (Username,Password) VALUES(?,?)")
	if err != nil {
		clog.Errorf("Error while preparing query: %v", err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(name, hash)
	if err != nil {
		clog.Errorf("Error inserting new data: %v", err)
	}

	var result string
	stmtOut, err := db.Prepare("SELECT 1 FROM Accounts WHERE Username = ?")
	err = stmtOut.QueryRow(cleanString(name)).Scan(&result)
	if err != nil {
		clog.Errorf("Issue while verifying data: %v", err)
		ui.MsgBoxError(loginwin, "Registration Failure", "The username you entered already exists.")
	} else {
		clog.Info("Account has been registered")
		ui.MsgBox(loginwin, "Registration Success", "You can now log in using your details.")
	}
}

func getPassword(name string) string {
	clog.Info("Trying to establish database connection")

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		clog.Fatalf("Could not establish database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		clog.Fatalf("Database connectivity issue: %v", err)
	}
	defer db.Close()

	var hash string
	stmtOut, err := db.Prepare("SELECT Password FROM Accounts WHERE Username = ?")
	if err != nil {
		clog.Errorf("Error while preparing query: %v", err)
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(cleanString(name)).Scan(&hash)
	if err != nil {
		clog.Errorf("Issue while scanning data: %v", err)
		ui.MsgBoxError(loginwin, "Login Failure", "The username you entered could not be found.")
	}

	return hash
}

func cleanString(input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		clog.Fatalf("Error while cleaning string: %v", err)
	}
	return reg.ReplaceAllString(input, "")
}

func controlUsername(entry *ui.Entry) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		clog.Fatalf("Error while cleaning string: %v", err)
	}
	if reg.MatchString(entry.Text()) {
		entry.SetText(cleanString(entry.Text()))
	}
}

func setupLoginForm() ui.Control {
	clog.Info("Creating the login form")
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	vbox.Append(entryForm, false)

	usernameIn := ui.NewEntry()
	usernameIn.OnChanged(controlUsername)
	passwordIn := ui.NewPasswordEntry()
	entryForm.Append("Username", usernameIn, true)
	entryForm.Append("Password", passwordIn, true)

	loginbtn := ui.NewButton("Login")
	loginbtn.OnClicked(func(*ui.Button) {
		if len(usernameIn.Text()) > 0 {
			clog.Info("Attempting to login user")
			hash := getPassword(usernameIn.Text())
			if len(hash) > 0 {
				err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordIn.Text()))
				if err != nil {
					clog.Errorf("Password authentication failed: %v", err)
					ui.MsgBoxError(loginwin, "Login Failure", "The password you entered is wrong.")
				} else {
					clog.Info("Successfully logged in!")
					ui.MsgBox(loginwin, "Login Success", "You have been successfully logged in.")
				}
			}
		}
	})

	vbox.Append(loginbtn, false)

	return vbox
}

func setupRegisterForm() ui.Control {
	clog.Info("Creating the registration form")
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	vbox.Append(entryForm, false)

	usernameIn := ui.NewEntry()
	usernameIn.OnChanged(controlUsername)
	entryForm.Append("Username", usernameIn, true)
	passwordIn := ui.NewPasswordEntry()
	entryForm.Append("Password", passwordIn, true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	registerbtn := ui.NewButton("Register")
	registerbtn.OnClicked(func(*ui.Button) {
		clog.Info("Registering user account")
		username := usernameIn.Text()
		hash := hashAndSalt([]byte(passwordIn.Text()))
		if len(username) == 0 || len(hash) == 0 {
			clog.Error("Either username or password is left empty")
			ui.MsgBoxError(loginwin, "Registration Failure", "Please enter an username and/or password.")
		} else {
			registerAccount(username, hash)
		}
	})
	hbox.Append(registerbtn, true)

	tosbtn := ui.NewButton("Terms of Service")
	tosbtn.OnClicked(func(*ui.Button) {
		clog.Info("Opening terms of service")
		// TODO: Create a proper ToS page and link it here
		openBrowser("https://netrve.net/")
	})
	hbox.Append(tosbtn, true)

	vbox.Append(hbox, false)

	return vbox
}

// SetupLogin() is the main function that setups the form and returns the window so it can be used in the main thread
func SetupLogin() *ui.Window {
	clog.Info("Preparing the login window")
	loginwin = ui.NewWindow(fmt.Sprintf("Loot Master v%s - Login", version.Version), 340, 220, true)
	loginwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		loginwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	loginwin.SetChild(tab)
	loginwin.SetMargined(true)

	tab.Append("Login", setupLoginForm())
	tab.SetMargined(0, true)

	tab.Append("Register", setupRegisterForm())
	tab.SetMargined(1, true)

	dataSource = database.GetDataSource()

	return loginwin
}
