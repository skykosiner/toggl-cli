package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func RemoveComments(jsonc []byte) []byte {
	re := regexp.MustCompile(`(?m)//.*$|/\*.*?\*/`)
	return re.ReplaceAll(jsonc, nil)
}

func MakeRequest(method, path, apiKey string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, "https://api.track.toggl.com/api/v9" + path, body)
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
