-- name: GetRepositoryByOwnerAndRepo :one
SELECT id, owner, repo, full_name, description, stars, forks, created_at, updated_at, last_fetched_at
FROM repositories
WHERE owner = $1 AND repo = $2;

-- name: UpsertRepository :exec
INSERT INTO repositories (owner, repo, full_name, description, stars, forks, created_at, updated_at, last_fetched_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
ON CONFLICT (owner, repo) 
DO UPDATE SET
    full_name = EXCLUDED.full_name,
    description = EXCLUDED.description,
    stars = EXCLUDED.stars,
    forks = EXCLUDED.forks,
    created_at = EXCLUDED.created_at,
    updated_at = NOW(),
    last_fetched_at = NOW();

-- name: GetAllRepositories :many
SELECT owner, repo, full_name, description, stars, forks, created_at, updated_at, last_fetched_at
FROM repositories
ORDER BY id;