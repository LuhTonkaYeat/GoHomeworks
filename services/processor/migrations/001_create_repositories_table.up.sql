CREATE TABLE IF NOT EXISTS repositories (
    id SERIAL PRIMARY KEY,
    owner VARCHAR(255) NOT NULL,
    repo VARCHAR(255) NOT NULL,
    full_name VARCHAR(512) NOT NULL,
    description TEXT,
    stars INT NOT NULL DEFAULT 0,
    forks INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_fetched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(owner, repo)
);

CREATE INDEX idx_repositories_owner_repo ON repositories(owner, repo);
CREATE INDEX idx_repositories_full_name ON repositories(full_name);