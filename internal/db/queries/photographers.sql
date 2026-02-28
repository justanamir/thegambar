-- name: ListPhotographers :many
SELECT id, name, specialty, city
FROM photographers
ORDER BY id;

-- name: GetPhotographer :one
SELECT id, name, specialty, city, bio, email, whatsapp, website
FROM photographers
WHERE id = $1;

-- name: InsertPhotographer :one
INSERT INTO photographers (name, specialty, city, bio, email, whatsapp, website)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;