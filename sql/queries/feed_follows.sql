-- name: CreateFeedFollow :one
WITH
  new_row as (
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
  users.name as user_name,
  feeds.name as feed_name
FROM
  users
  INNER JOIN new_row ON new_row.user_id = users.id
  INNER JOIN feeds on new_row.feed_id = feeds.id;
