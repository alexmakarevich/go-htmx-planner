-- EVENTS:

-- name: GetCalendaEvent :one
SELECT * FROM calendar_events
WHERE id = ? LIMIT 1;

-- name: ListCalendaEvents :many
SELECT * FROM calendar_events
ORDER BY date_time;

-- name: CreateCalendaEvent :one
INSERT INTO calendar_events (
  title, date_time, owner_id
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateCalendaEvent :exec
UPDATE calendar_events
set title = ?,
date_time = ?,
owner_id = ?
WHERE id = ?;

-- name: DeleteCalendaEvent :exec
DELETE FROM calendar_events
WHERE id = ?;

-- USERS:

-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: FindUser :one
SELECT * FROM users
WHERE user_name = ? AND password = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY user_name;

-- name: CreateUser :one
INSERT INTO users (
  user_name, password
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
set user_name = ?,
password = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- SESSIONS:

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = ? LIMIT 1;

-- name: GetSessionWithUser :one
SELECT sessions.id as session_id, sessions.user_id, users.user_name  FROM sessions INNER JOIN users ON sessions.user_id = users.id
WHERE sessions.id = ? LIMIT 1;

-- name: ListSessions :many
SELECT * FROM sessions
ORDER BY id;

-- name: CreateSession :one
INSERT INTO sessions (
  id, user_id
) VALUES (
  ?, ?
)
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;