-- name: CreateRssFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetAllRssFeeds :many
SELECT * FROM feeds;


-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at DESC 
LIMIT $1;


-- name: MarkFeedFetched :one
UPDATE feeds
SET updated_at = $2, last_fetched_at = $3
WHERE id = $1
RETURNING *;
