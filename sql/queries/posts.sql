-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING *;
-- 

-- name: GetPostsForUser :many
SELECT p.*, f.name AS feed_name FROM posts p
JOIN feed_follows ff ON ff.feed_id = p.feed_id
JOIN feeds f ON p.feed_id = f.id
WHERE ff.user_id = $1
ORDER BY p.published_at DESC
LIMIT $2;
--
