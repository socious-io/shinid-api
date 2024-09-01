INSERT INTO media(user_id, url, filename)
VALUES($1, $2, $3)
RETURNING *