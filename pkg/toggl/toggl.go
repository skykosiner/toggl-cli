package toggl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/skykosiner/toggl-cli/pkg/utils"
)

type SavedTimer struct {
	Name       string   `json:"name"`
	ProjectID  int      `json:"project_id"`
	Tags       []string `json:"tags"`
	Description string   `json:"description"`
}

type Toggl struct {
	ApiKey      string       `json:"api_key"`
	WorkspaceID int          `json:"workspace_id"`
	SavedTimers []SavedTimer `json:"saved_timers"`
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

func NewToggl() (Toggl, error) {
	var t Toggl
	configDir, err := os.UserConfigDir()
	if err != nil {
		return t, err
	}

	configPath := path.Join(configDir, "toggl", "config.jsonc")
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(utils.RemoveComments(bytes), &t)
	return t, err
}

func (t Toggl) UpdateConfigFile() {
	bytes, err := json.MarshalIndent(t, " ", "  ")
	if err != nil {
		slog.Error("Couldn't update config file", "erorr", err)
		os.Exit(1)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Couldn't update config file", "erorr", err)
		os.Exit(1)
	}

	configPath := path.Join(configDir, "toggl", "config.jsonc")

	if err := os.WriteFile(configPath, bytes, 0644); err != nil {
		slog.Error("Couldn't update config file", "erorr", err)
		os.Exit(1)
	}
}

func (t Toggl) selectTags() []string {
	var tags []struct{
		Name string `json:"name"`
	}
	var tagSlice []string

	resp, err := utils.MakeRequest(http.MethodGet, fmt.Sprintf("/workspaces/%d/tags", t.WorkspaceID), t.ApiKey, nil)
	if err != nil {
		return tagSlice
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return tagSlice
	}

	if err := json.Unmarshal(bytes, &tags); err != nil {
		fmt.Println(err)
		return tagSlice
	}

	idxs, err := fuzzyfinder.FindMulti(tags,
	func(i int) string {
		return fmt.Sprintf("%s", tags[i].Name)
	})

	for _, idx := range idxs {
		tagSlice = append(tagSlice, tags[idx].Name)
	}

	return tagSlice
}

func (t Toggl) newEntry(body NewTimeEntry) error {
	b, _ := json.Marshal(body)
	resp, err := utils.MakeRequest(http.MethodPost, fmt.Sprintf("/workspaces/%d/time_entries", t.WorkspaceID), t.ApiKey, bytes.NewBuffer(b))
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

func (t Toggl) ResumeEntry() error {
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

	if err := t.newEntry(NewTimeEntry{
		CreatedWith: "toggl cli",
		Description: cachedTimer.Description,
		Duration:    -1,
		ProjectID:   cachedTimer.ProjectID,
		Start:       time.Now().UTC().Format(time.RFC3339),
		Tags:        cachedTimer.Tags,
		WorkspaceID: t.WorkspaceID,
	}); err != nil {
		return err
	}

	return nil
}

func (t Toggl) StartSaved() error {
	idx, err := fuzzyfinder.Find(t.SavedTimers,
		func(i int) string {
			return fmt.Sprintf("Name: %s ProjectID: %d, Tags: %s, Description: %s", t.SavedTimers[i].Name, t.SavedTimers[i].ProjectID, strings.Join(t.SavedTimers[i].Tags, ", "), t.SavedTimers[i].Description)
		})

	if err != nil {
		if err.Error() == "abort" {
			return nil
		}

		return err
	}

	if err := t.newEntry(NewTimeEntry{
		CreatedWith: "toggl cli",
		Description: t.SavedTimers[idx].Description,
		Duration:    -1,
		ProjectID:   t.SavedTimers[idx].ProjectID,
		Start:       time.Now().UTC().Format(time.RFC3339),
		Tags:        t.SavedTimers[idx].Tags,
		WorkspaceID: t.WorkspaceID,
	}); err != nil {
		return err
	}

	return nil
}

func (t Toggl) Start() error {
	var description string

	projects := t.fetchProjects()
	idx, err := fuzzyfinder.Find(projects,
	func(i int) string {
		return fmt.Sprintf("%s", projects[i].Name)
	})

	if err != nil {
		if err.Error() == "abort" {
			return nil
		}

		return err
	}

	tagSlice := []string{}
	if utils.AskInput("Would you like to add tags? (y/n)") == "y" {
		tagSlice = t.selectTags()
	}

	description = utils.AskInput("Description? (leave blank for no description")

	if err := t.newEntry(NewTimeEntry{
		CreatedWith: "toggl cli",
		Description: description,
		Duration:    -1,
		ProjectID:   projects[idx].ID,
		Start:       time.Now().UTC().Format(time.RFC3339),
		Tags:        tagSlice,
		WorkspaceID: t.WorkspaceID,
	}); err != nil {
		return err
	}

	return nil
}
