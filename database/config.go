package database

import (
	"encoding/json"
	"io/ioutil"

	"github.com/dakraid/LooM/clog"
)

type Config struct {
	Login    string `json:"login"`
	Pass     string `json:"pass"`
	IP       string `json:"ip"`
	Protocol string `json:"protocol"`
	Database string `json:"database"`
}

var jsonIn Config

const config = "database.json"

func ReadJson() {
	dat, err := ioutil.ReadFile(config)
	if err != nil {
		clog.Fatalf("Couldn't load database config: %v",err)
	}

	err = json.Unmarshal(dat, &jsonIn)
	if err != nil {
		clog.Fatalf("Error while loading config: %v",err)
	}
}

func GetDataSource() string {
	dataSource := jsonIn.Login + ":" + jsonIn.Pass + "@" + jsonIn.Protocol + "(" + jsonIn.IP + ")/" + jsonIn.Database
	return dataSource
}