SELECT u.*,
m.url as "avatar.url",
m.filename "avatar.filename"
FROM users u
LEFT JOIN media m ON u.avatar_id=m.id
WHERE username = $1