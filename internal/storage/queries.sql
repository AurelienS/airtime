-- name: GetUserWithGoogleId :one
SELECT
  *
FROM users
WHERE
  google_id = $1
LIMIT 1;

-- name: UpsertUser :one
INSERT INTO users (google_id, email, NAME, picture_url)
VALUES
  ($1, $2, $3, $4) ON CONFLICT (google_id) DO
UPDATE
SET
  email = EXCLUDED.email,
  NAME = EXCLUDED.name,
  picture_url = EXCLUDED.picture_url,
  updated_at = NOW() RETURNING *;

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

-- name: InsertFlight :execresult
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
  ) RETURNING id;

-- name: InsertGlider :exec
INSERT INTO gliders(NAME, user_id)
VALUES
  ($1, $2);

-- name: UpdateDefaultGlider :exec
UPDATE
  users
SET
  default_glider_id = $1
WHERE
  id = $2;

-- name: InsertFlightStats :one
INSERT INTO flight_statistics (
    total_thermic_time,
    total_flight_time,
    max_climb,
    max_climb_rate,
    total_climb,
    average_climb_rate,
    number_of_thermals,
    percentage_thermic,
    max_altitude
  )
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
  ) RETURNING id;
