package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func MakeRequest(method, path, apiKey string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, "https://api.track.toggl.com/api/v9"+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}

func AskInput(msg string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(msg)
	input, err := reader.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}

	input = input[:len(input)-1]
	return input
}

func getTerminalSize() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		// Fallback to default size
		return 24
	}

	var rows int
	fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &rows)
	return rows
}

func PrintBanner(msg string) {
	cols := getTerminalSize()
	x := (cols-len(msg))/2

	fmt.Printf("\033[%dH\033[34m%s\033[0m\n", x, msg)
}
