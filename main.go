package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/skykosiner/toggl-cli/pkg/toggl"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Short: "toggl - toggl cli",
		Use: "toggl",
	}

	t, err := toggl.NewToggl()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	commands := []cobra.Command{
		{
			Use: "status",
			Short: "Get the curent tracking status",
			Run: func(cmd *cobra.Command, args []string) {
				curr := t.GetCurrentTimer()
				fmt.Printf("%s %s %s %s\n", curr.Description, strings.Join(curr.Tags, ", "), curr.GetProjectName(t.ApiKey, t.WorkspaceID), curr.GetDuration())
			},
		},
		{
			Use: "pause",
			Short: "Pause the current entry",
			Run: func(cmd *cobra.Command, args []string) {
				curr := t.GetCurrentTimer()
				if err := curr.Stop(t.ApiKey, t.WorkspaceID, true); err != nil {
					fmt.Println(err)
					return
				}
			},
		},
		{
			Use: "stop",
			Short: "Stop the current entry",
			Run: func(cmd *cobra.Command, args []string) {
				curr := t.GetCurrentTimer()
				if err := curr.Stop(t.ApiKey, t.WorkspaceID, false); err != nil {
					fmt.Println(err)
					return
				}
			},
		},
		{
			Use: "resume",
			Short: "Resume the paused time entry",
			Run: func(cmd *cobra.Command, args []string) {
				if err := t.ResumeEntry(); err != nil{
					fmt.Println(err)
					return
				}
			},
		},
		{
			Use: "start-saved",
			Short: "Start new time entry from your saved timers",
			Run: func(cmd *cobra.Command, args []string) {
				if err := t.StartSaved(); err != nil{
					fmt.Println(err)
					return
				}
			},
		},
		{
			Use: "start",
			Short: "Start new time entry",
			Run: func(cmd *cobra.Command, args []string) {
				if err := t.Start(); err != nil{
					fmt.Println(err)
					return
				}
			},
		},
		{
			Use: "new-saved",
			Short: "Save a new time entry",
			Run: func(cmd *cobra.Command, args []string) {
				if err := t.NewSaveTimer(); err != nil{
					fmt.Println(err)
					return
				}
			},
		},
	}

	for _, command := range commands {
		rootCmd.AddCommand(&command)
	}

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}
