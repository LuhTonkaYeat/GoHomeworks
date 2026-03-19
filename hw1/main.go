package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stargazers  int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error! Use: go run main.go owner/repo")
		return
	}

	path := strings.Split(os.Args[1], "/")
	if len(path) != 2 {
		fmt.Println("Error! Use: owner/repo")
		return
	}

	owner := path[0]
	repo := path[1]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	fmt.Printf("Requesting: %s\n", url)

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error while requesting: %v\n", err)
		return
	}
	defer response.Body.Close()

	fmt.Printf("Response status: %s\n", response.Status)

	switch response.StatusCode {
	case http.StatusOK:
		fmt.Println("Repository found! Parsing data...")

	case http.StatusNotFound:
		message := fmt.Sprintf(`Repository not found!
The repository '%s/%s' does not exist or is private.

Possible issues:
- Check the spelling (GitHub is case-sensitive!)
- Your repo should be: LuhTonkaYeat/GoHW1 (not luhTONKAyeat/gohw1)
- Make sure the repository is public`, owner, repo)

		fmt.Println(message)
		return

	default:
		fmt.Printf("Error! GitHub returned status: %d\n", response.StatusCode)
		fmt.Println("Please try again later.")
		return
	}

	var repoInfo Repository
	err = json.NewDecoder(response.Body).Decode(&repoInfo)
	if err != nil {
		fmt.Printf("Error while parsing JSON: %v\n", err)
		return
	}

	description := repoInfo.Description
	if description == "" {
		description = "No description"
	}

	createdAt := repoInfo.CreatedAt

	result := fmt.Sprintf(`
===== Repository Information =====
Name: %s
Description: %s
Stars: %d
Forks: %d
Created: %s`,
		repoInfo.Name,
		description,
		repoInfo.Stargazers,
		repoInfo.Forks,
		createdAt,
	)

	fmt.Println(result)
}
