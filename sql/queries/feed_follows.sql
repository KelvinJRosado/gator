-- name: CreateFeedFollow :one
WITH
  new_row AS (
    INSERT INTO
      feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES
      ($1, $2, $3, $4, $5)
    RETURNING
      *
  )
SELECT
  new_row.id,
  new_row.created_at,
  new_row.updated_at,
  users.name AS user_name,
  feeds.name AS feed_name
FROM
  users
  INNER JOIN new_row ON new_row.user_id = users.id
  INNER JOIN feeds ON new_row.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT
  feed_follows.id,
  feed_follows.created_at,
  feed_follows.updated_at,
  users.name AS user_name,
  feeds.name AS feed_name,
  feeds.url AS feed_url
FROM
  users
  INNER JOIN feed_follows ON feed_follows.user_id = users.id
  INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
  users.name = $1;

-- name: DeleteFeedFollowsForUser :exec
DELETE FROM feed_follows
WHERE
  user_id = $1
  AND feed_id = $2;
