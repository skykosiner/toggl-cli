package main

import (
	"encoding/json"
	"os"
	"path"
	"regexp"
)

func removeComments(jsonc []byte) []byte {
	re := regexp.MustCompile(`(?m)//.*$|/\*.*?\*/`)
	return re.ReplaceAll(jsonc, nil)
}

type SavedTimer struct {
	Name        string   `json:"name"`
	ProjectID   int`json:"project_id"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

type Config struct {
	ApiKey      string       `json:"api_key"`
	WorkspaceID int          `json:"workspace_id"`
	SavedTimers []SavedTimer `json:"saved_timers"`
}

func NewConfig() (Config, error) {
	var config Config
	configDir, err := os.UserConfigDir()
	if err != nil {
		return config, err
	}

	configPath := path.Join(configDir, "toggl", "config.jsonc")
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(removeComments(bytes), &config)
	return config, err
}
