package database

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dakraid/LooM/clog"
)

type config struct {
	Login    string `json:"login"`
	Pass     string `json:"pass"`
	IP       string `json:"ip"`
	Protocol string `json:"protocol"`
	Database string `json:"database"`
}

var jsonIn config

const configfile = "database.json"

func readJson() {
	dat, err := ioutil.ReadFile(configfile)
	if err != nil {
		clog.Fatalf("Couldn't load database config: %v", err)
	}

	err = json.Unmarshal(dat, &jsonIn)
	if err != nil {
		clog.Fatalf("Error while loading config: %v", err)
	}
}

// GetDataSource returns a string to be used with the MySQL driver to establish a connection
func GetDataSource() string {
	readJson()
	dataSource := jsonIn.Login + ":" + jsonIn.Pass + "@" + jsonIn.Protocol + "(" + jsonIn.IP + ")/" + jsonIn.Database
	return dataSource
}
