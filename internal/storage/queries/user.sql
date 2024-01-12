-- name: GetUserWithGoogleId :one
SELECT
  *
FROM users
WHERE
  google_id = $1
LIMIT 1;

-- name: UpsertUser :one
INSERT INTO users (
    google_id,
    email,
    NAME,
    picture_url
  )
VALUES
  ($1, $2, $3, $4) ON CONFLICT (google_id) DO
UPDATE
SET
  email = EXCLUDED.email,
  NAME = EXCLUDED.name,
  picture_url = EXCLUDED.picture_url RETURNING *;