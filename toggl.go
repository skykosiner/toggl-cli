package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CurrentEntry struct {
	ProjectName string `json:"project_name"`
	Tags []string  `json:"tags"`
	Start string `json:"start"`
	Description string `json:"description"`
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