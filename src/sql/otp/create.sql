INSERT INTO otps(user_id, code)
VALUES ($1, $2)
RETURNING *