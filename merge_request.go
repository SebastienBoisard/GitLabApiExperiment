package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"encoding/json"
)

// MergeRequest contains the Merge Request data
// Used https://mholt.github.io/json-to-go/
type MergeRequest struct {
	ID           int       `json:"id"`
	Iid          int       `json:"iid"`
	ProjectID    int       `json:"project_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TargetBranch string    `json:"target_branch"`
	SourceBranch string    `json:"source_branch"`
	Upvotes      int       `json:"upvotes"`
	Downvotes    int       `json:"downvotes"`
	Author       struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		ID        int    `json:"id"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"author"`
	Assignee                 interface{}   `json:"assignee"`
	SourceProjectID          int           `json:"source_project_id"`
	TargetProjectID          int           `json:"target_project_id"`
	Labels                   []interface{} `json:"labels"`
	WorkInProgress           bool          `json:"work_in_progress"`
	Milestone                interface{}   `json:"milestone"`
	MergeWhenBuildSucceeds   bool          `json:"merge_when_build_succeeds"`
	MergeStatus              string        `json:"merge_status"`
	Sha                      string        `json:"sha"`
	MergeCommitSha           string        `json:"merge_commit_sha"`
	Subscribed               bool          `json:"subscribed"`
	UserNotesCount           int           `json:"user_notes_count"`
	ApprovalsBeforeMerge     interface{}   `json:"approvals_before_merge"`
	ShouldRemoveSourceBranch interface{}   `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch  bool          `json:"force_remove_source_branch"`
	WebURL                   string        `json:"web_url"`
}

/**
 * getMergedRequests retrieves the Merge Requests that have been already merged.
 *
 * Doc: https://docs.gitlab.com/ee/api/merge_requests.html#list-merge-requests
 */
func getMergedRequests(gitlabToken string, gitlabUrl string, projectName string) ([]MergeRequest, error) {

	projectName = url.QueryEscape(projectName)

	restUrl := fmt.Sprintf("%s/api/v3/projects/%s/merge_requests?state=merged", gitlabUrl, projectName)

	// Build the request
	req, err := http.NewRequest("GET", restUrl, nil)
	if err != nil {
		log.Println("NewRequest: ", err)
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", gitlabToken)

	// Create a HTTP Client for control over HTTP client headers, redirect policy, and other settings.
	client := &http.Client{}

	// Send an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
		return nil, err
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var mergedRequests []MergeRequest

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&mergedRequests); err != nil {
		log.Println(err)
		return nil, err
	}

	return mergedRequests, nil
}
