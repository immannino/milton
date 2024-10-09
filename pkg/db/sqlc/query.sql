-- name: GetWebsites :many
SELECT * FROM website;

-- name: UpsertWebsite :one
INSERT INTO website (url) VALUES (?) 
ON CONFLICT (url) DO UPDATE SET
 updated = CURRENT_TIMESTAMP
RETURNING *;

-- name: UpsertTag :one
INSERT INTO tag (website_id, tag_type, name, property, value)
VALUES (?, ?, ?, ?, ?)
RETURNING *;