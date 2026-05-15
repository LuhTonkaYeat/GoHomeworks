package postgres

import (
	"context"
	
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/repository"
)

type Repository struct {
	queries *repository.Queries
	pool    *pgxpool.Pool
}

func NewRepository(connString string) (*Repository, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	
	queries := repository.New(pool)
	
	return &Repository{
		queries: queries,
		pool:    pool,
	}, nil
}

func (r *Repository) Close() {
	r.pool.Close()
}

func (r *Repository) GetQueries() *repository.Queries {
	return r.queries
}