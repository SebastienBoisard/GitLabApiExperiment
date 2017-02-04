package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"encoding/json"
)

// Commit contains the commit data
type Commit struct {
	ID          string    `json:"id"`
	ShortID     string    `json:"short_id"`
	Title       string    `json:"title"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	CreatedAt   time.Time `json:"created_at"`
	Message     string    `json:"message"`
}

/**
 * getCommits retrieves all the commits of a repository branch.
 *
 * Doc: https://docs.gitlab.com/ee/api/commits.html#list-repository-commits
 */
func getCommits(gitlabToken string, gitlabUrl string, projectName string, commitName string) ([]Commit, error) {

	projectName = url.QueryEscape(projectName)

	commitName = url.QueryEscape(commitName)

	restUrl := fmt.Sprintf("%s/api/v3/projects/%s/repository/commits?ref_name=%s",
		gitlabUrl, projectName, commitName)

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
	var commits []Commit

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		log.Println(err)
		return nil, err
	}

	return commits, nil
}
