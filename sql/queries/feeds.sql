-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
		gen_random_uuid(),
		NOW(),
		NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedbyUrl :one
select * from feeds where url = $1;

-- name: CreateFeedFollow :one
with inserted_feed_follows as (
	insert into feed_follows (id, created_at, updated_at, user_id, feed_id)
	values (
		gen_random_uuid(),
		NOW(),
		NOW(),
		$1,
		$2
	)
	returning *
)
select inserted_feed_follows.*,
		feeds.name as feed_name,
		users.name as user_name
	from inserted_feed_follows
	inner join users on inserted_feed_follows.user_id = users.id
	inner join feeds on inserted_feed_follows.feed_id = feeds.id;

-- name: GetFeedsByUser :many
select ff.*, 
	u.name as user_name,
	f.name as feed_name 
from feed_follows ff
inner join users u on u.id = ff.user_id
inner join feeds f on f.id = ff.feed_id
where ff.user_id = $1;

-- name: DeleteFollowByUser :exec
delete from feed_follows 
where user_id = $1 and feed_id = $2;
