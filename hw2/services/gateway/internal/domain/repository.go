package domain

import "time"

type Repository struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stars       int       `json:"stars"`
	Forks       int       `json:"forks"`
	CreatedAt   time.Time `json:"created_at"`
}
