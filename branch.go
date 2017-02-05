package GitLabApiExperiment

import (
	"encoding/json"
	"fmt"
 	"log"
	"net/http"
	"net/url"
	"time"
)

// Branch contains the branch data
type Branch struct {
	Name   string `json:"name"`
	Commit struct {
		ID             string    `json:"id"`
		Message        string    `json:"message"`
		ParentIds      []string  `json:"parent_ids"`
		AuthoredDate   time.Time `json:"authored_date"`
		AuthorName     string    `json:"author_name"`
		AuthorEmail    string    `json:"author_email"`
		CommittedDate  time.Time `json:"committed_date"`
		CommitterName  string    `json:"committer_name"`
		CommitterEmail string    `json:"committer_email"`
	} `json:"commit"`
	Merged             bool `json:"merged"`
	Protected          bool `json:"protected"`
	DevelopersCanPush  bool `json:"developers_can_push"`
	DevelopersCanMerge bool `json:"developers_can_merge"`
}

/**
 * getBranches retrieves branches from a project.
 *
 * Doc: https://docs.gitlab.com/ee/api/branches.html#list-repository-branches
 */
func GetBranches(gitlabToken string, gitlabUrl string, projectName string) ([]Branch, error) {

	projectName = url.QueryEscape(projectName)

	restUrl := fmt.Sprintf("%s/api/v3/projects/%s/repository/branches", gitlabUrl, projectName)

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
	var branches []Branch

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		log.Println(err)
		return nil, err
	}

	return branches, nil
}
