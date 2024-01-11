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


-- interval cast because it was seen as a bigint ?
-- name: GetTotalFlightTime :one
SELECT sum(total_flight_time)::text
FROM   flights fl
JOIN   flight_statistics fls ON fls.id = fl.flight_statistics_id
where fl.user_id = $1;
