SELECT o.*,
m.url as "logo.url",
m.filename "logo.filename",
(SELECT status FROM kyb_verifications kv WHERE o.id = kv.organization_id ORDER BY created_at DESC LIMIT 1) AS verification_status
FROM organizations o
LEFT JOIN media m ON o.logo_id=m.id
WHERE o.id IN(?)