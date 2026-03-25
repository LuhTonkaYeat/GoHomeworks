package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LuhTonkaYeat/GoHomeworks/hw2/internal/collector/domain"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.github.com",
	}
}

type githubRepoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stargazers  int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func (c *Client) FetchRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("repository %s/%s not found", owner, repo)
		case http.StatusForbidden:
			return nil, fmt.Errorf("rate limit exceeded or access denied")
		default:
			return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
		}
	}

	var githubResp githubRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&githubResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, githubResp.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	return &domain.Repository{
		Name:        githubResp.Name,
		Description: githubResp.Description,
		Stars:       githubResp.Stargazers,
		Forks:       githubResp.Forks,
		CreatedAt:   createdAt,
	}, nil
}
