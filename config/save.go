package config

import (
	"encoding/json"
	"log"
	"os"
)

// Save save current status and config to file
func Save() {
	SaveStatus()
	SaveConfig()
}

// SaveStatus save and indent Status to status.json
func SaveStatus() {
	b, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		log.Println("Failed to save status...")
		return
	}
	os.WriteFile("status.json", b, 0755)
}

// SaveConfig save and indent Config to config.json
func SaveConfig() {
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Println("Failed to save config...")
		return
	}
	os.WriteFile("config.json", b, 0755)
}
