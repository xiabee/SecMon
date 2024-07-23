package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	githubToken   = "your_github_token" // GitHub Personal Access Token
	repoOwner     = "repo_owner"        // 仓库所有者
	repoName      = "repo_name"         // 仓库名称
	checkInterval = 3600 * time.Second  // 检查间隔时间
)

var (
	labels   = []string{"security", "vulnerability"}
	keywords = []string{"security", "vulnerability", "漏洞", "安全"}
)

// Issue represents a GitHub issue
type Issue struct {
	Title string `json:"title"`
	URL   string `json:"html_url"`
}

// getIssues fetches the issues from the GitHub repository
func getIssues() ([]Issue, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?labels=%s&state=open",
		repoOwner, repoName, strings.Join(labels, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+githubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch issues: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var issues []Issue
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

// filterIssues filters the issues based on keywords in the title
func filterIssues(issues []Issue) []Issue {
	var filtered []Issue
	for _, issue := range issues {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(issue.Title), strings.ToLower(keyword)) {
				filtered = append(filtered, issue)
				break
			}
		}
	}
	return filtered
}

// notify prints the filtered issues to the console
func notify(issues []Issue) {
	for _, issue := range issues {
		fmt.Printf("New security issue detected: %s - %s\n", issue.Title, issue.URL)
	}
}

func main() {
	for {
		issues, err := getIssues()
		if err != nil {
			fmt.Printf("Error fetching issues: %v\n", err)
			time.Sleep(checkInterval)
			continue
		}

		filteredIssues := filterIssues(issues)
		if len(filteredIssues) > 0 {
			notify(filteredIssues)
		}

		time.Sleep(checkInterval)
	}
}
