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
