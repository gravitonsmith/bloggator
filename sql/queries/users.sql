-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: DeleteUser :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;
