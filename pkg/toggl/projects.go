package toggl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/skykosiner/toggl-cli/pkg/utils"
)

type Project struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func (t Toggl) fetchProjects() []Project {
	var projects []Project
	resp, err := utils.MakeRequest(http.MethodGet, fmt.Sprintf("/workspaces/%d/projects", t.WorkspaceID), t.ApiKey, nil)
	if err != nil {
		return projects
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return projects
	}

	if err := json.Unmarshal(bytes, &projects); err != nil {
		return projects
	}

	projects = slices.DeleteFunc(projects, func(item Project) bool {
		return !item.Active
	})


	return projects
}
