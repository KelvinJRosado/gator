-- name: CreatePost :one
INSERT INTO
  posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
  *;

-- name: GetPostsForUser :many
SELECT
  posts.id,
  posts.created_at,
  posts.updated_at,
  posts.title,
  posts.url,
  posts.description,
  posts.published_at,
  posts.feed_id,
  users.name AS user_name,
  feeds.name AS feed_name,
  feeds.url AS feed_url,
  feeds.id AS feed_id
FROM
  users
  INNER JOIN feeds ON users.id = feeds.user_id
  INNER JOIN posts ON posts.feed_id = feeds.id
WHERE
  users.name = $1
ORDER BY
  posts.published_at DESC
LIMIT
  $2;
