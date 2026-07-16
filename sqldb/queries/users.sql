-- name: IsOwnerPresent :one
SELECT COUNT(*) FROM users WHERE is_owner = 1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
	email, first_name, last_name,
	password, is_owner, added_by
)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
	email      = ?,
	first_name = ?,
	last_name  = ?,
	avatar     = ?,
	about_me   = ?
WHERE id = ?
RETURNING *;

-- name: UpdatePassword :one
UPDATE users
SET password = ?
WHERE id = ?
RETURNING *;

-- name: TransferOwnership :exec
UPDATE users SET is_owner = CASE
    WHEN id = ? THEN 0
    WHEN id = ? THEN 1
END
WHERE id IN (?, ?);
