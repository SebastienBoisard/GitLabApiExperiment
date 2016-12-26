package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
	"time"
	"flag"
)

// Used https://mholt.github.io/json-to-go/
// MergeRequest contains the Merge Request data
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
	Author struct {
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

// Branch contains the branch data
type Branch struct {
	Name string `json:"name"`
	Commit struct {
		ID             string `json:"id"`
		Message        string `json:"message"`
		ParentIds      []string `json:"parent_ids"`
		AuthoredDate   time.Time `json:"authored_date"`
		AuthorName     string `json:"author_name"`
		AuthorEmail    string `json:"author_email"`
		CommittedDate  time.Time `json:"committed_date"`
		CommitterName  string `json:"committer_name"`
		CommitterEmail string `json:"committer_email"`
	} `json:"commit"`
	Protected          bool `json:"protected"`
	DevelopersCanPush  bool `json:"developers_can_push"`
	DevelopersCanMerge bool `json:"developers_can_merge"`
}

// Commit contains the commit data
type Commit struct {
	ID          string `json:"id"`
	ShortID     string `json:"short_id"`
	Title       string `json:"title"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	CreatedAt   time.Time `json:"created_at"`
	Message     string `json:"message"`
}

/**
 * getMergedRequests retrieves the Merge Requests that have been already merged.
 *
 * Doc: https://docs.gitlab.com/ee/api/merge_requests.html#list-merge-requests
 */
func getMergedRequests(gitlabToken string, gitlabUrl string, projectName string) (error, []MergeRequest) {

	projectName = url.QueryEscape(projectName)

	restUrl := fmt.Sprintf("%s/api/v3/projects/%s/merge_requests?state=merged&private_token=%s", gitlabUrl, projectName, gitlabToken)

	// Build the request
	req, err := http.NewRequest("GET", restUrl, nil)
	if err != nil {
		log.Println("NewRequest: ", err)
		return err, nil
	}

	// Create a HTTP Client for control over HTTP client headers, redirect policy, and other settings.
	client := &http.Client{}

	// Send an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
		return err, nil
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var mergedRequests []MergeRequest

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&mergedRequests); err != nil {
		log.Println(err)
		return err, nil
	}

	return nil, mergedRequests
}

/**
 * getBranches retrieves branches from a project.
 *
 * Doc: https://docs.gitlab.com/ee/api/branches.html#list-repository-branches
 */
func getBranches(gitlabToken string, gitlabUrl string, projectName string) (error, []Branch) {

	projectName = url.QueryEscape(projectName)

	restUrl := fmt.Sprintf("%s/api/v3/projects/%s/repository/branches?private_token=%s", gitlabUrl, projectName, gitlabToken)

	// Build the request
	req, err := http.NewRequest("GET", restUrl, nil)
	if err != nil {
		log.Println("NewRequest: ", err)
		return err, nil
	}

	// Create a HTTP Client for control over HTTP client headers, redirect policy, and other settings.
	client := &http.Client{}

	// Send an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
		return err, nil
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var branches []Branch

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		log.Println(err)
		return err, nil
	}

	return nil, branches
}

/**
 * getCommits retrieves all the commits of a repository branch.
 *
 * Doc: https://docs.gitlab.com/ee/api/commits.html#list-repository-commits
 */
func getCommits(gitlabToken string, gitlabUrl string, projectName string, commitName string) (error, []Commit) {

	projectName = url.QueryEscape(projectName)

	commitName = url.QueryEscape(commitName)

	restUrl := fmt.Sprintf("%s/api/v3/projects/%s/repository/commits?ref_name=%s&private_token=%s",
		gitlabUrl, projectName, commitName, gitlabToken)

	// Build the request
	req, err := http.NewRequest("GET", restUrl, nil)
	if err != nil {
		log.Println("NewRequest: ", err)
		return err, nil
	}

	// Create a HTTP Client for control over HTTP client headers, redirect policy, and other settings.
	client := &http.Client{}

	// Send an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
		return err, nil
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	var commits []Commit

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		log.Println(err)
		return err, nil
	}

	return nil, commits
}



func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("Error: no configuration file not found")
		return
	}

	gitlabToken := viper.GetString("gitlab.token")
	gitlabUrl := viper.GetString("gitlab.url")
	projectName := viper.GetString("project.name")

	var commandName string
	flag.StringVar(&commandName, "command_name", "", "the command to execute: merged_requests, all_branches, commits")

	flag.Parse()

	if commandName == "" {
		fmt.Println("Error: command_name is missing")
		fmt.Println("")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	switch commandName {
	case "merged_requests":
		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		// Get all the merged requests from a GitLab project
		err, mergedRequests := getMergedRequests(gitlabToken, gitlabUrl, projectName)
		if err != nil {
			log.Println("Error: can't get the merged requests [", err, "]")
			return
		}

		for _, r := range mergedRequests {
			fmt.Println("merged requests title = ", r.Title)
			fmt.Println("                status = ", r.State)
			fmt.Println("                created at = ", r.CreatedAt)
			fmt.Println("                source branch = ", r.SourceBranch)
			fmt.Println("                target branch = ", r.TargetBranch)
		}

	case "all_branches":
		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		// Get all the branches from a GitLab project
		err, branches := getBranches(gitlabToken, gitlabUrl, projectName)
		if err != nil {
			log.Println("Error: can't get the branches [", err, "]")
			return
		}

		for _, r := range branches {
			fmt.Println("branch name = ", r.Name)
		}

	case "commits":
		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		// Get all the commits from a specific branch of the GitLab project
		err, commits := getCommits(gitlabToken, gitlabUrl, projectName, "cert-fast-load")
		if err != nil {
			log.Println("Error: can't get the commits [", err, "]")
			return
		}

		for _, r := range commits {
			fmt.Printf("commit date = %s  title = %s  \n", r.CreatedAt, r.Title)
		}
		
	default:
		fmt.Println("Error: unknown command name")
		fmt.Println("")
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}


/*
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Get which branch was merged in which branch
	for _, branch := range branches {
		for _, mergedRequest := range mergedRequests {
			fmt.Printf("compare %s vs %s\n", branch.Name, mergedRequest.Title)
			if branch.Name == mergedRequest.SourceBranch {
				fmt.Printf("branch '%s' was merged into branch '%s' on %s\n", branch.Name, mergedRequest.TargetBranch,
					mergedRequest.UpdatedAt.Format("2006-01-02 15:04"))
			}
		}
	}
*/
}
