UPDATE users
SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    bio = COALESCE($4, bio),
    job_title = COALESCE($5, job_title),
    phone = COALESCE($6, phone),
    username = COALESCE($7, username),
    avatar_id = COALESCE($8, avatar_id)
WHERE id = $1
RETURNING *;
