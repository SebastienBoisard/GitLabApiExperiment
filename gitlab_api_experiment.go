package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
	"time"
)

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


type Branch struct {
   Name string `json:"name"`
   Commit struct {
      ID string `json:"id"`
      Message string `json:"message"`
      ParentIds []string `json:"parent_ids"`
      AuthoredDate time.Time `json:"authored_date"`
      AuthorName string `json:"author_name"`
      AuthorEmail string `json:"author_email"`
      CommittedDate time.Time `json:"committed_date"`
      CommitterName string `json:"committer_name"`
      CommitterEmail string `json:"committer_email"`
   } `json:"commit"`
   Protected bool `json:"protected"`
   DevelopersCanPush bool `json:"developers_can_push"`
   DevelopersCanMerge bool `json:"developers_can_merge"`
}

func getMergedRequests(gitlabToken string, projectName string) (error, []MergeRequest) {

   projectName = url.QueryEscape(projectName)

   url := fmt.Sprintf("http://www.gitlab.com/api/v3/projects/%s/merge_requests?state=merged&private_token=%s", projectName, gitlabToken)

   // Build the request
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      log.Fatal("NewRequest: ", err)
      return err, nil
   }

   // Create a HTTP Client for control over HTTP client headers, redirect policy, and other settings.
   client := &http.Client{}

   // Send an HTTP request and returns an HTTP response
   resp, err := client.Do(req)
   if err != nil {
      log.Fatal("Do: ", err)
      return err, nil
   }

   // Defer the closing of the body
   defer resp.Body.Close()

   // Fill the record with the data from the JSON
   var record []MergeRequest

   // Use json.Decode for reading streams of JSON data
   if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
      log.Println(err)
      return err, nil
   }

   return nil, record
}

func getBranches(gitlabToken string, projectName string) (error, []Branch) {

   projectName = url.QueryEscape(projectName)

   url := fmt.Sprintf("http://www.gitlab.com/api/v3/projects/%s/repository/branches?private_token=%s", projectName, gitlabToken)

   // Build the request
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      log.Fatal("NewRequest: ", err)
      return err, nil
   }

   // Create a HTTP Client for control over HTTP client headers, redirect policy, and other settings.
   client := &http.Client{}

   // Send an HTTP request and returns an HTTP response
   resp, err := client.Do(req)
   if err != nil {
      log.Fatal("Do: ", err)
      return err, nil
   }

   // Defer the closing of the body
   defer resp.Body.Close()

   // Fill the record with the data from the JSON
   var record []Branch

   // Use json.Decode for reading streams of JSON data
   if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
      log.Println(err)
      return err, nil
   }

   return nil, record
}

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("Error: no configuration file not found")
		return
	}

	gitlabToken := viper.GetString("connection.token")


   err, mergedRequests := getMergedRequests(gitlabToken, "technomancy/bussard")
   if err != nil {
      log.Println("Error: can't get the merged requests [", err, "]")
      return      
   }

   for _, r := range mergedRequests {
      fmt.Println("merged requests title = ", r.Title)
   }


   err, branches := getBranches(gitlabToken, "gnutls/gnutls")
   if err != nil {
      log.Println("Error: can't get the branches [", err, "]")
      return      
   }

   for _, r := range branches {
      fmt.Println("branch name = ", r.Name)
   }
}
