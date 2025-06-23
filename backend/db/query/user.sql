-- name: CreateUser :exec
INSERT INTO users (id, username, password_hash, email, roles)
VALUES (?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT id, username, email, roles
FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT id, username, email, roles
FROM users
ORDER BY username;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: UpdateUsername :exec
UPDATE users SET username = ? WHERE id = ?;

-- name: UpdateEmail :exec
UPDATE users SET email = ? WHERE id = ?;

-- name: UpdateRoles :exec
UPDATE users SET roles = ? WHERE id = ?;

-- name: GetUserByEmailWithoutPassword :one
SELECT id, username, email, roles
FROM users
WHERE email = ?;

-- name: GetUserByEmailWithPassword :one
SELECT id, username, email, password_hash, roles
FROM users
WHERE email = ?;
