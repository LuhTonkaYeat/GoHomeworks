package domain

import "time"

type Repository struct {
	Name        string
	Description string
	Stars       int
	Forks       int
	CreatedAt   time.Time
}
