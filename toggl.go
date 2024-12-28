package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type CurrentEntry struct {
	ID int `json:"id"`
	ProjectID int `json:"project_id"`
	Tags []string  `json:"tags"`
	Start string `json:"start"`
	Description string `json:"description"`
}

func (c CurrentEntry) cache() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		slog.Error("Error getting cache dir", "error", err)
		return
	}

	togglCache := filepath.Join(cacheDir, "toggl")
	if err := os.MkdirAll(togglCache, 0755); err != nil {
		slog.Error("Error creating cache dir.", "error", err, "path", togglCache)
		return
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		slog.Error("Error marshalling JSON of the timer.", "error", err, "timer", c)
		return
	}

	if err := os.WriteFile(filepath.Join(togglCache, "toggl.json"), bytes, 0644); err != nil {
		slog.Error("Error updating JSON cache file.", "error", err, "timer", c)
		return
	}
}

func (c CurrentEntry) GetDuration() string {
	parsedTime, err := time.Parse(time.RFC3339, c.Start)
	if err != nil {
		return ""
	}

	currentTime := time.Now()
	duration := currentTime.Sub(parsedTime)

	// Extract hours and minutes
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60


	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func (c CurrentEntry) GetProjectName(apiKey string, workspaceID int) string {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/projects/%d", workspaceID, c.ProjectID), nil)
	if err != nil {
		slog.Error("Error making request", "error", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
	}


	var r struct{
		Name string `json:"name"`
	}
	if err := json.Unmarshal(bytes, &r); err != nil {
		slog.Error("Error making request", "error", err)
	}

	return r.Name
}

func GetCurrentEntry(apiKey string) (CurrentEntry, error) {
	var curr CurrentEntry
	req, err := http.NewRequest(http.MethodGet, "https://api.track.toggl.com/api/v9/me/time_entries/current", nil)
	if err != nil {
		return curr, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return curr, err
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return curr, err
	}

	err = json.Unmarshal(bytes, &curr)
	return curr, err
}
