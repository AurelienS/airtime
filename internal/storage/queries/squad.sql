-- name: InsertSquad :one
INSERT INTO squads (NAME)
VALUES
    ($1) RETURNING id;

-- name: InsertSquadMember :exec
INSERT INTO squad_members (squad_id, user_id, admin, joined_at)
VALUES
    ($1, $2, $3, NOW());

-- name: FindAllSquadForUser :many
SELECT
    sqlc.embed(s),
    sqlc.embed(sm)
FROM squads s
    JOIN squad_members sm ON s.id = sm.squad_id
WHERE
    sm.user_id = $1;