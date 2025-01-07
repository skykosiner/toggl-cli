package toggl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skykosiner/toggl-cli/pkg/utils"
)

type NewTag struct {
	Name        string `json:"name"`
	WorkspaceID int    `json:"workspace_id"`
}

func (t *Toggl) NewTags(tags []string) error {
	for _, tag := range tags {
		newTag := NewTag{
			Name:        tag,
			WorkspaceID: t.WorkspaceID,
		}
		jB, _ := json.Marshal(newTag)

		resp, err := utils.MakeRequest(http.MethodPost, fmt.Sprintf("/workspaces/%d/tags", t.WorkspaceID), t.ApiKey, bytes.NewBuffer(jB))
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Something went wrong got status code %d", resp.StatusCode)
		}

		fmt.Printf("Created tag %s\n", tag)
	}

	return nil
}
