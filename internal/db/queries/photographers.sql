-- name: ListPhotographers :many
SELECT id, name, specialty, city, cover_url
FROM photographers
ORDER BY id;

-- name: GetPhotographer :one
SELECT id, name, specialty, city, bio, email, whatsapp, website, avatar_url, cover_url
FROM photographers
WHERE id = $1;

-- name: InsertPhotographer :one
INSERT INTO photographers (name, specialty, city, bio, email, whatsapp, website, edit_token)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdatePhotographerPhotos :one
UPDATE photographers
SET avatar_url = $2, cover_url = $3
WHERE id = $1
RETURNING id;

-- name: GetPhotographerByToken :one
SELECT id, name, specialty, city, bio, email, whatsapp, website, avatar_url, cover_url, edit_token
FROM photographers
WHERE edit_token = $1;

-- name: UpdatePhotographer :one
UPDATE photographers
SET name = $2, specialty = $3, city = $4, bio = $5, email = $6, whatsapp = $7, website = $8
WHERE edit_token = $1
RETURNING id;