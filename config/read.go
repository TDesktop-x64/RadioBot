package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
)

var (
	status Status
	config Config
	cron   = gocron.NewScheduler(time.UTC)
)

func Read() {
	ReadConfig()
	ReadStatus()
}

func ReadConfig() {
	if err := LoadConfig(); err != nil {
		log.Fatal(err)
	}
}

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

func ReadStatus() {
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
