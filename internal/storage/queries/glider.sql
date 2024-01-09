-- name: GetGliders :many
SELECT
  *
FROM gliders
WHERE
  user_id = $1;

-- name: InsertGlider :exec
INSERT INTO gliders(NAME, user_id)
VALUES
  ($1, $2);
