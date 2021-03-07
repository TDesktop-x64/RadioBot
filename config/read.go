package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/c0re100/RadioBot/utils"
	"github.com/go-co-op/gocron"
)

var (
	status Status
	config Config
	cron   = gocron.NewScheduler(time.UTC)
)

// Read read config and status JSON
func Read() {
	readConfig()
	readStatus()
}

func readConfig() {
	if err := LoadConfig(); err != nil {
		log.Fatal(err)
	}
	// Check port is valid
	port := GetWebPort()
	utils.CheckPortIsValid("Web Server", port)
	port = GetBeefWebPort()
	utils.CheckPortIsValid("Beefweb", port)
	// Check setting is valid
	compareUpdateVoteTime()
}

// LoadConfig load config.json to Config
func LoadConfig() error {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}
	e := json.Unmarshal(b, &config)
	if e != nil {
		return e
	}
	return nil
}

func initStatus() {
	b := []byte("{}")
	os.WriteFile("status.json", b, 0755)
	json.Unmarshal(b, &status)
}

func readStatus() {
	if b, err := ioutil.ReadFile("status.json"); err == nil {
		e := json.Unmarshal(b, &status)
		if e != nil {
			log.Println("status.json is broken...resetting")
			initStatus()
		}
	} else {
		log.Println("status.json not found...initializating")
		initStatus()
	}
	cron.Every(1).Minute().Do(func() {
		SaveStatus()
	})
	cron.StartAsync()
}
