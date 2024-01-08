-- name: GetUserWithGoogleId :one
SELECT
  *
FROM users
WHERE
  google_id = $1
LIMIT 1;

-- name: UpsertUser :exec
INSERT INTO users (google_id, email, NAME, picture_url)
VALUES
  ($1, $2, $3, $4) ON CONFLICT (google_id) DO
UPDATE
SET
  email = EXCLUDED.email,
  NAME = EXCLUDED.name,
  picture_url = EXCLUDED.picture_url,
  updated_at = NOW();

-- name: GetFlights :many
SELECT
  *
FROM flights
WHERE
  user_id = $1;

-- name: GetGliders :many
SELECT
  *
FROM gliders
WHERE
  user_id = $1;

-- name: InsertFlight :exec
INSERT INTO flights (
    DATE,
    takeoff_location,
    igc_file_path,
    user_id,
    glider_id,
    flight_statistics_id
  )
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
  );

-- name: InsertGlider :exec
INSERT INTO gliders(NAME, user_id)
VALUES
  ($1, $2);

-- name: UpdateDefaultGlider :exec
UPDATE users
SET default_glider_id = $1
WHERE id = $2;
