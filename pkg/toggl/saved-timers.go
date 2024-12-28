package toggl

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/skykosiner/toggl-cli/pkg/utils"
)

func (t *Toggl) NewSavedTimer() error {
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

func (t *Toggl) DeleteSavedTimer() error {
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

	t.SavedTimers = slices.Delete(t.SavedTimers, idx, idx+1)

	t.UpdateConfigFile()
	return nil
}
