-- name: CreateFeed :one
INSERT INTO
  feeds (id, created_at, updated_at, name, url, user_id)
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING
  *;

-- name: GetAllFeeds :many
SELECT
  feeds.id,
  feeds.created_at,
  feeds.updated_at,
  feeds.name,
  feeds.url,
  users.name AS user_name
FROM
  feeds
  INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT
  feeds.id,
  feeds.created_at,
  feeds.updated_at,
  feeds.name,
  feeds.url,
  users.name AS user_name
FROM
  feeds
  INNER JOIN users ON feeds.user_id = users.id
WHERE
  feeds.url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET
  updated_at = $1,
  last_fetched_at = $1
WHERE
  id = $2
RETURNING
  *;

-- name: GetNextFeedToFetch :one
SELECT
  *
FROM
  feeds
ORDER BY
  last_fetched_at ASC NULLS FIRST
LIMIT
  1;
