package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	config, err := NewConfig()
	if err != nil {
	    slog.Error("Sorry there was an error processing your config. Check the README.", "error", err)
		return
	}

	rootCmd := &cobra.Command{
		Short: "toggl - toggl cli",
		Use: "toggl",
	}

	commands := []cobra.Command{
		{
			Use: "status",
			Short: "Get the curent tracking status",
			Run: func(cmd *cobra.Command, args []string) {
				curr, err := GetCurrentEntry(config.ApiKey)
				if err != nil {
					slog.Error("Error getting the current status.", "error", err)
					return
				}

				if curr.Description != "" {
					fmt.Printf("%s: %s: %s\n", curr.GetProjectName(config.ApiKey, config.WorkspaceID), curr.Description, curr.GetDuration())
				} else {
					fmt.Printf("%s: %s\n", curr.GetProjectName(config.ApiKey, config.WorkspaceID), curr.GetDuration())
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
