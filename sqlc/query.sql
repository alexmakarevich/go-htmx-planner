-- EVENTS:

-- name: GetCalendarEvent :one
SELECT * FROM calendar_events
WHERE id = ? LIMIT 1;

-- name: GetCalendarEventWithOwner :one
SELECT calendar_events.*, users.user_name as owner_name  FROM calendar_events LEFT JOIN users ON calendar_events.owner_id = users.id
WHERE calendar_events.id = ? LIMIT 1;

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

-- TODO: separate passwords cleanly

-- USERS:

-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- -- name: ListUsersInRelationToThisEvent :many
-- SELECT users.*, filtered_participations.event_id as event_id FROM users LEFT JOIN 
-- (
--   SELECT * from participations
--   WHERE event_id = ?
-- ) as filtered_participations
-- ON filtered_participations.user_id = users.id;

-- name: ListUsersInRelationToThisEvent :many
SELECT users.*, participations.event_id as event_id, participations.status as status FROM users LEFT JOIN participations
ON participations.user_id = users.id AND participations.event_id = ?;

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


-- PARTICIPATIONS

-- name: AddParticipant :one
INSERT INTO participations (
  user_id, event_id, status
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- FYI: sqlite upserts are broken in sqlc

-- name: UpdateParticipant :exec
UPDATE participations SET status = @status
WHERE user_id = @user_id AND event_id = @event_id;

-- name: InviteParticipants :many
UPDATE participations SET status = 'invited'
WHERE event_id = @event_id
RETURNING *;

-- name: DeleteParticipant :exec
DELETE FROM participations
WHERE user_id = ? AND event_id = ?;

-- name: GetParticipantsByEventId :many
SELECT users.* FROM participations INNER JOIN users ON participations.user_id = users.id
WHERE participations.event_id = ? AND participations.status = ?;

-- name: SearchUsersExcludingParticipants :many
SELECT users.* FROM users LEFT JOIN 
( SELECT * FROM participations WHERE event_id = CAST(@event_id AS INTEGER)) parts
ON parts.user_id = users.id
WHERE parts.user_id IS NULL AND user_name LIKE CONCAT('%', CAST(@query AS TEXT), '%')
LIMIT 10
; -- the cast is there for sqlc

-- name: SearchParticipants :many
SELECT users.* FROM users
WHERE user_name LIKE CONCAT('%', CAST(@query AS TEXT), '%'); -- the cast is there for sqlc
