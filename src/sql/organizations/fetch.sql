SELECT o.*,
m.url as "logo.url",
m.filename "logo.filename"
FROM organizations o
LEFT JOIN media m ON o.logo_id=m.id
WHERE o.id IN(?)