ALTER TABLE photographers
ADD COLUMN edit_token TEXT NOT NULL DEFAULT '';

UPDATE photographers
SET edit_token = md5(random()::text || id::text)
WHERE edit_token = '';

ALTER TABLE photographers
ADD CONSTRAINT photographers_edit_token_unique UNIQUE (edit_token);