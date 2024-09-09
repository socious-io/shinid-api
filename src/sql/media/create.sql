INSERT INTO media(user_id, url, filename)
VALUES($1, $2, $3)
ON CONFLICT (url)
DO UPDATE SET filename = $3
RETURNING *