-- name: GetUserWithGoogleId :one
SELECT
  *
FROM users
WHERE
  google_id = $1
LIMIT 1;

-- name: UpsertUser :one
INSERT INTO users (google_id, email, NAME, picture_url, default_glider_id)
VALUES
  ($1, $2, $3, $4, $5) ON CONFLICT (google_id) DO
UPDATE
SET
  email = EXCLUDED.email,
  NAME = EXCLUDED.name,
  picture_url = EXCLUDED.picture_url,
  default_glider_id = $5,
  updated_at = NOW() RETURNING *;

-- name: UpdateDefaultGlider :exec
UPDATE
  users
SET
  default_glider_id = $1
WHERE
  id = $2;
