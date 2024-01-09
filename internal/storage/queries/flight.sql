-- name: GetFlights :many
SELECT
  *
FROM flights
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
