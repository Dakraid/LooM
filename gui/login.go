package gui

import (
	"database/sql"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/google/logger"
	"golang.org/x/crypto/bcrypt"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	_ "github.com/go-sql-driver/mysql"

	"github.com/dakraid/LooM/database"
	"github.com/dakraid/LooM/version"
)

var (
	loginwin    *ui.Window
	dataSource  string
)

func openBrowser(url string) {
	logger.Info("Attempting to open link in browser")
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
		logger.Fatalf("Error while opening URL: %v", err)
	}
}

func hashAndSalt(pwd []byte) string {
	logger.Info("Generating hash and salt for the password")
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)

	if err != nil {
		logger.Fatalf("Error while generating hash: %v", err)
	}

	return string(hash)
}

// TODO: Handle the database connection in its own function or own package
func connectDatabase() *sql.DB {
	logger.Info("Trying to establish database connection")

	db, err := sql.Open("mysql",dataSource)
	if err != nil {
		logger.Fatalf("Could not establish database connection: %v",err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatalf("Database connectivity issue: %v",err)
	}
	defer db.Close()

	return db
}

func registerAccount(name, hash string) {
	if len(name) == 0 || len(hash) == 0 {
		logger.Error("Either username or password is left empty")
		ui.MsgBoxError(loginwin,"Registration Failure","Please enter an username and/or password.")
	}
	logger.Info("Trying to establish database connection")

	db, err := sql.Open("mysql",dataSource)
	if err != nil {
		logger.Fatalf("Could not establish database connection: %v",err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatalf("Database connectivity issue: %v",err)
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO Accounts (Username,Password) VALUES(?,?)")
	if err != nil {
		logger.Errorf("Error while preparing query: %v",err)
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(name, hash)
	if err != nil {
		logger.Errorf("Error inserting new data: %v",err)
	}

	var result string
	stmtOut, err := db.Prepare("SELECT 1 FROM Accounts WHERE Username = ?")
	err = stmtOut.QueryRow(cleanString(name)).Scan(&result)
	if err != nil {
		logger.Errorf("Issue while verifying data: %v",err)
		ui.MsgBoxError(loginwin,"Registration Failure","The username you entered already exists.")
	} else {
		logger.Info("Account has been registered")
		ui.MsgBox(loginwin,"Registration Success","You can now log in using your details.")
	}
}

func getPassword(name string) string {
	logger.Info("Trying to establish database connection")

	db, err := sql.Open("mysql",dataSource)
	if err != nil {
		logger.Fatalf("Could not establish database connection: %v",err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatalf("Database connectivity issue: %v",err)
	}
	defer db.Close()

	var hash string
	stmtOut, err := db.Prepare("SELECT Password FROM Accounts WHERE Username = ?")
	if err != nil {
		logger.Errorf("Error while preparing query: %v",err)
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(cleanString(name)).Scan(&hash)
	if err != nil {
		logger.Errorf("Issue while scanning data: %v",err)
		ui.MsgBoxError(loginwin,"Login Failure","The username you entered could not be found.")
	}

	return hash
}

func cleanString(input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		logger.Fatalf("Error while cleaning string: %v",err)
	}
	return reg.ReplaceAllString(input, "")
}

func controlUsername(entry *ui.Entry) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		logger.Fatalf("Error while cleaning string: %v",err)
	}
	if reg.MatchString(entry.Text()) {
		entry.SetText(cleanString(entry.Text()))
	}
}

func setupLoginForm() ui.Control {
	logger.Info("Creating the login form")
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
			logger.Info("Attempting to login user")
			hash := getPassword(usernameIn.Text())
				if len(hash) > 0 {
				err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordIn.Text()))
				if err != nil {
					logger.Errorf("Password authentication failed: %v",err)
					ui.MsgBoxError(loginwin,"Login Failure","The password you entered is wrong.")
				} else {
					logger.Info("Successfully logged in!")
					ui.MsgBox(loginwin,"Login Success","You have been successfully logged in.")
				}
			}
		}
	})

	vbox.Append(loginbtn,false)

	return vbox
}

func setupRegisterForm() ui.Control {
	logger.Info("Creating the registration form")
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
		logger.Info("Registering user account")
		username := usernameIn.Text()
		hash := hashAndSalt([]byte(passwordIn.Text()))
		registerAccount(username,hash)
	})
	hbox.Append(registerbtn,true)


	tosbtn := ui.NewButton("Terms of Service")
	tosbtn.OnClicked(func(*ui.Button) {
		logger.Info("Opening terms of service")
		// TODO: Create a proper ToS page and link it here
		openBrowser("https://netrve.net/")
	})
	hbox.Append(tosbtn,true)

	vbox.Append(hbox,false)


	return vbox
}

func setupLogin() {
	logger.Info("Preparing the login window")
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

	loginwin.Show()
}

func ShowLogin() {
	ui.Main(setupLogin)
}