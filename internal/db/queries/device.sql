-- name: GetDevices :many
SELECT id, name, status FROM devices;

-- name: InsertDevice :one
INSERT INTO devices (name, status) VALUES ($1, $2) RETURNING *;
