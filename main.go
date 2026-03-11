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

	response, error := http.Get(url)
	if error != nil {
		fmt.Printf("Error while requesting: %v\n", error)
		return
	}
	defer response.Body.Close()

	fmt.Printf("Response status: %s\n", response.Status)

	switch response.StatusCode {
	case 200:
		fmt.Println("Repository found! Parsing data...")

	case 404:
		fmt.Println("Repository not found!")
		fmt.Printf("The repository '%s/%s' does not exist or is private.\n", owner, repo)
		fmt.Println("\nPossible issues:")
		fmt.Println("- Check the spelling (GitHub is case-sensitive!)")
		fmt.Printf("- Your repo should be: LuhTonkaYeat/GoHW1 (not luhTONKAyeat/gohw1)\n")
		fmt.Println("- Make sure the repository is public")
		return

	default:
		fmt.Printf("Error! GitHub returned status: %d\n", response.StatusCode)
		fmt.Println("Please try again later.")
		return

	}

	var repoInfo Repository
	err := json.NewDecoder(response.Body).Decode(&repoInfo)
	if err != nil {
		fmt.Printf("Error while parsing JSON: %v\n", err)
		return
	}

	fmt.Println("\n===== Repository Information =====")
	fmt.Printf("Name: %s\n", repoInfo.Name)
	if repoInfo.Description != "" {
		fmt.Printf("Description: %s\n", repoInfo.Description)
	} else {
		fmt.Printf("Description: No description\n")
	}
	fmt.Printf("Stars: %d\n", repoInfo.Stargazers)
	fmt.Printf("Forks: %d\n", repoInfo.Forks)
	fmt.Printf("Created: %s\n", repoInfo.CreatedAt)
}
