package main

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

func main() {
	config, err := NewConfig()
	if err != nil {
	    slog.Error("Sorry there was an error processing your config. Check the README.", "error", err)
		return
	}

	cobra

	curr, err := GetCurrentEntry(config.ApiKey)
	if err != nil {
	    slog.Error("Sorry there was an error getting the current time entry.", "error", err)
		return
	}

	fmt.Println(curr.GetDuration())
}
