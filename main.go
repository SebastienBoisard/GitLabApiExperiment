package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

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
		mergedRequests, err := getMergedRequests(gitlabToken, gitlabUrl, projectName)
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
		branches, err := getBranches(gitlabToken, gitlabUrl, projectName)
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
		commits, err := getCommits(gitlabToken, gitlabUrl, projectName, "cert-fast-load")
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
}
