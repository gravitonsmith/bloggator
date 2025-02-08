-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: DeleteUser :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateFeed :one
INSERT INTO feed (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
		$5,
		$6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feed.*, users.name AS user_name
FROM feed
JOIN users ON feed.user_id = users.id;
