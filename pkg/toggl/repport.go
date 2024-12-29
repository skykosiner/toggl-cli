package toggl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/skykosiner/toggl-cli/pkg/utils"
)

type ReportBody struct {
	EndDate   string `json:"end_date"`
	StartDate string `json:"start_date"`
}

type ReportType int

const (
	Daily ReportType = iota
	Week
	Monthly
	Yearly
)

type Report struct {
	ProjectID       int `json:"project_id"`
	TrackedSeconds int `json:"tracked_seconds"`
}

type ReportFinal struct {
	Project string
	Hours, Minutes, Seconds int
}

func (r Report) getProjectName(apiKey string, workspaceID int) string {
	resp, err := utils.MakeRequest(http.MethodGet, fmt.Sprintf("/workspaces/%d/projects/%d", workspaceID, r.ProjectID), apiKey, nil)
	if err != nil {
		slog.Error("Error when making request to get project name", "error", err)
		return ""
	}

	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	var re struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(bytes, &re); err != nil {
		return "No Project"
	}

	return re.Name
}

func (t Toggl) GetReport(r ReportType) {
	var startDate string
	switch r {
	case Daily:
		startDate = time.Now().Format("2006-01-02")
		utils.PrintBanner(`
╺━╸╺━╸╺━╸   ╺┳┓┏━┓╻ ╻   ╺━╸╺━╸╺━╸
╺━╸╺━╸╺━╸    ┃┃┣━┫┗┳┛   ╺━╸╺━╸╺━╸
╺━╸╺━╸╺━╸   ╺┻┛╹ ╹ ╹    ╺━╸╺━╸╺━╸
`)
	case Week:
		startDate = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		utils.PrintBanner(`
╺━╸╺━╸╺━╸   ╻ ╻┏━╸┏━╸╻┏    ╺━╸╺━╸╺━╸
╺━╸╺━╸╺━╸   ┃╻┃┣╸ ┣╸ ┣┻┓   ╺━╸╺━╸╺━╸
╺━╸╺━╸╺━╸   ┗┻┛┗━╸┗━╸╹ ╹   ╺━╸╺━╸╺━╸
`)
	case Monthly:
		startDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
		utils.PrintBanner(`
╺━╸╺━╸╺━╸   ┏┳┓┏━┓┏┓╻╺┳╸╻ ╻   ╺━╸╺━╸╺━╸
╺━╸╺━╸╺━╸   ┃┃┃┃ ┃┃┗┫ ┃ ┣━┫   ╺━╸╺━╸╺━╸
╺━╸╺━╸╺━╸   ╹ ╹┗━┛╹ ╹ ╹ ╹ ╹   ╺━╸╺━╸╺━╸

`)
	case Yearly:
		startDate = time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
		utils.PrintBanner(`
╺━╸╺━╸╺━╸    ╻ ╻┏━╸┏━┓┏━┓    ╺━╸╺━╸╺━╸
╺━╸╺━╸╺━╸    ┗┳┛┣╸ ┣━┫┣┳┛    ╺━╸╺━╸╺━╸
╺━╸╺━╸╺━╸     ╹ ┗━╸╹ ╹╹┗╸    ╺━╸╺━╸╺━╸
`)
	default:
		fmt.Println("Invalid ReportType")
		return
	}

	rp := ReportBody{
		StartDate: startDate,
		EndDate:   time.Now().Format("2006-01-02"),
	}

	jB, _ := json.Marshal(rp)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.track.toggl.com/reports/api/v3/workspace/%d/projects/summary", t.WorkspaceID), bytes.NewBuffer(jB))
	if err != nil {
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(t.ApiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		os.Exit(1)
	}

	defer resp.Body.Close()

	var reports []Report
	bytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bytes, &reports)

	var reportFinals []ReportFinal
	for _, report := range reports {
		duration := time.Duration(report.TrackedSeconds) * time.Second

		// Extract hours, minutes, and seconds
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		seconds := int(duration.Seconds()) % 60

		reportFinals = append(reportFinals, ReportFinal{
			Project: report.getProjectName(t.ApiKey, t.WorkspaceID),
			Hours: hours,
			Minutes: minutes,
			Seconds: seconds,
		})
	}

	sort.Slice(reportFinals, func(i, j int) bool {
		if reportFinals[i].Hours != reportFinals[j].Hours {
			return reportFinals[i].Hours > reportFinals[j].Hours
		}
		if reportFinals[i].Minutes != reportFinals[j].Minutes {
			return reportFinals[i].Minutes > reportFinals[j].Minutes
		}
		return reportFinals[i].Seconds > reportFinals[j].Seconds
	})

	for _, report := range reportFinals {
		fmt.Printf("%s: %02d:%02d:%02d\n", report.Project, report.Hours, report.Minutes, report.Seconds)
	}
}
