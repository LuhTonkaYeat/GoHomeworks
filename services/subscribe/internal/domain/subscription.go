package domain

import "time"

type Subscription struct {
	ID        int64
	UserID    string
	Owner     string
	Repo      string
	CreatedAt time.Time
}
