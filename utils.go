package main

import (
	"io"
	"net/http"
)

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
