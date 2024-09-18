INSERT INTO recipients 
(first_name, last_name, email, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *