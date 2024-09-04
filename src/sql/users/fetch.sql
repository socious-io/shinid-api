SELECT u.*,
m.url as "avatar.url",
m.filename "avatar.filename"
FROM users
LEFT JOIN media m ON u.avatar_id=m.id
WHERE id IN (?)