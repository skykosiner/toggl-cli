package main

import (
	"fmt"
	"log/slog"
)

func main() {
	config, err := NewConfig()
	if err != nil {
	    slog.Error("Sorry there was an error processing your config. Check the README.", "error", err)
	}

	fmt.Println(config)
}
