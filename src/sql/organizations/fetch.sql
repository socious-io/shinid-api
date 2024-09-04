SELECT o.*,
m.url as "logo.url",
m.filename "logo.filename"
FROM organizations
LEFT JOIN media m ON o.logo_id=m.id
WHERE id IN(?)