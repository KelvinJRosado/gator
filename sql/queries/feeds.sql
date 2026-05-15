-- name: CreateFeed :one
INSERT INTO
  feeds (id, created_at, updated_at, name, url, user_id)
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING
  *;

-- name: GetAllFeeds :many
SELECT
  /* sql-formatter-disable */
  sqlc.embed(feeds),
  sqlc.embed(users)
/* sql-formatter-enable */
FROM
  feeds
  INNER JOIN users ON feeds.user_id = users.id;
