package config

import (
	"encoding/json"
	"log"
	"os"
)

func Save() {
	SaveStatus()
	SaveConfig()
}

func SaveStatus() {
	b, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		log.Println("Failed to save status...")
		return
	}
	os.WriteFile("status.json", b, 0755)
}

func SaveConfig() {
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Println("Failed to save config...")
		return
	}
	os.WriteFile("config.json", b, 0755)
}
