-- name: GetLink :one
SELECT * FROM links WHERE id = $1;

-- name: GetLinkByOriginalUrl :one
SELECT * FROM links WHERE original_url = $1;

-- name: CreateLink :one
INSERT INTO links (id, original_url) VALUES ($1, $2)
RETURNING *;

-- name: DeleteLink :exec
DELETE FROM links WHERE id = $1;