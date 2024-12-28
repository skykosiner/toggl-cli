package main

import (
	"bytes"
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
	ID          int      `json:"id"`
	ProjectID   int      `json:"project_id"`
	Tags        []string `json:"tags"`
	Start       string   `json:"start"`
	Description string   `json:"description"`
}

type NewTimeEntry struct {
	CreatedWith string   `json:"created_with"`
	Description string   `json:"description"`
	Duration    int      `json:"duration"`
	ProjectID   int      `json:"project_id"`
	Start       string   `json:"start"`
	Tags        []string `json:"tags"`
	WorkspaceID int      `json:"workspace_id"`
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
	resp, err := MakeRequest(http.MethodGet, fmt.Sprintf("/workspaces/%d/projects/%d", workspaceID, c.ProjectID), apiKey, nil)
	if err != nil {
		slog.Error("Error when making request to get project name", "error", err)
		return ""
	}

	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	var r struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(bytes, &r); err != nil {
		slog.Error("Error making request", "error", err)
	}

	return r.Name
}

func (c CurrentEntry) Stop(apiKey string, workspaceID int, pause bool) error {
	if pause {
		c.cache()
	}

	resp, err := MakeRequest(http.MethodPatch, fmt.Sprintf("/workspaces/%d/time_entries/%dv/stop", workspaceID, c.ID), apiKey, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Couldn't stop timer please try again.")
	}

	return nil
}

func GetCurrentEntry(apiKey string) (CurrentEntry, error) {
	var curr CurrentEntry
	resp, err := MakeRequest(http.MethodPatch, "/me/time_entries/current", apiKey, nil)
	if err != nil {
		return curr, err
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return curr, err
	}

	if string(bytes) == "null" {
		// This means the timer isn't running when the bytes are null, so just
		// exit
		os.Exit(1)
	}

	err = json.Unmarshal(bytes, &curr)
	return curr, err
}

func NewEntry(body NewTimeEntry, apiKey string, workspaceID int) error {
	b, _ := json.Marshal(body)
	resp, err := MakeRequest(http.MethodPost, fmt.Sprintf("/workspaces/%d/time_entries", workspaceID), apiKey, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	fmt.Println(string(bytes))

	if resp.StatusCode != 200 {
		return errors.New("Couldn't start new time entry")
	}

	return nil
}

func ResumeEntry(apiKey string, workspaceID int) error {
	// First let's check if a cache exists, if it doesn't we can't resume the entry
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	togglCache := filepath.Join(cacheDir, "toggl", "toggl.json")
	bytes, err := os.ReadFile(togglCache)
	if err != nil {
		return err
	}

	var cachedTimer CurrentEntry
	if err := json.Unmarshal(bytes, &cachedTimer); err != nil {
		return err
	}

	if err := NewEntry(NewTimeEntry{
		CreatedWith: "toggl cli",
		Description: cachedTimer.Description,
		Duration:    -1,
		ProjectID:   cachedTimer.ProjectID,
		Start:       time.Now().UTC().Format(time.RFC3339),
		Tags:        cachedTimer.Tags,
		WorkspaceID: workspaceID,
	}, apiKey, workspaceID); err != nil {
		return err
	}

	return nil
}
