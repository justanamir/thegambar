-- name: ListPhotographers :many
SELECT id, name, specialty, city, cover_url
FROM photographers
ORDER BY id;

-- name: GetPhotographer :one
SELECT id, name, specialty, city, bio, email, whatsapp, website, avatar_url, cover_url
FROM photographers
WHERE id = $1;

-- name: InsertPhotographer :one
INSERT INTO photographers (name, specialty, city, bio, email, whatsapp, website)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdatePhotographerPhotos :one
UPDATE photographers
SET avatar_url = $2, cover_url = $3
WHERE id = $1
RETURNING id;