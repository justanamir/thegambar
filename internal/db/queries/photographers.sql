-- name: ListPhotographers :many
SELECT id, name, specialty, city
FROM photographers
ORDER BY id;

-- name: GetPhotographer :one
SELECT id, name, specialty, city, bio, email, whatsapp, website
FROM photographers
WHERE id = $1;