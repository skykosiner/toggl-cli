package toggl

import (
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/skykosiner/toggl-cli/pkg/utils"
)

func (t *Toggl) NewSaveTimer() error {
	name := utils.AskInput("Enter the name of the saved timer")
	if len(name) == 0 {
		fmt.Println("Please provide a name for the saved timer.")
		return nil
	}

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

	tags := t.selectTags()
	description := utils.AskInput("Description? (leave blank for no description")

	t.SavedTimers = append(t.SavedTimers, SavedTimer{
		Name: name,
		ProjectID: projects[idx].ID,
		Tags: tags,
		Description: description,
	})

	t.UpdateConfigFile()
	return nil
}
