package main

import (
	"net/http"
)

func MakeRequest(method, path, apiKey string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.track.toggl.com/api/v9" + path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}
