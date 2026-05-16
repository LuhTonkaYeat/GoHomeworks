package domain

import "time"

type Repository struct {
	Name        string
	Description string
	Stars       int
	Forks       int
	CreatedAt   time.Time
}

type DBRepository struct {
	ID            int64
	Owner         string
	Repo          string
	FullName      string
	Description   string
	Stars         int
	Forks         int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastFetchedAt time.Time
}