package dto

type RepoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
