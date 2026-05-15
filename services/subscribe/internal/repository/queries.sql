-- name: CreateSubscription :exec
INSERT INTO subscriptions (user_id, owner, repo)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, owner, repo) DO NOTHING;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE user_id = $1 AND owner = $2 AND repo = $3;

-- name: GetSubscriptions :many
SELECT owner, repo FROM subscriptions
WHERE user_id = $1;

-- name: SubscriptionExists :one
SELECT EXISTS (
    SELECT 1 FROM subscriptions
    WHERE user_id = $1 AND owner = $2 AND repo = $3
);