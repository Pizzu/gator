-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *
)
SELECT
    inserted_feed_follow.*,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow
INNER JOIN users u on user_id = u.id
INNER JOIN feeds f on feed_id = f.id;

-- name: GetFeedFollowsForUser :many
SELECT fs.*, f.name AS feed_name, u.name AS user_name
FROM feed_follows fs
INNER JOIN users u on u.id = fs.user_id
INNER JOIN feeds f on f.id = fs.feed_id
WHERE fs.user_id = $1;