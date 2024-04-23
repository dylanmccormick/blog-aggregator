-- name: CreateRssFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;


-- name: GetAllRssFeeds :many
SELECT * FROM feeds;
