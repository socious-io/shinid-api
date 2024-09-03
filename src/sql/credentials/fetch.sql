SELECT
  cv.*,
  row_to_json(u.*) AS created,
  row_to_json(o.*) AS organization,
  row_to_json(cs.*) AS schema,
  row_to_json(r.*) AS recipient
FROM credentials cv 
LEFT JOIN users u ON u.id = cv.created_id
LEFT JOIN organizations o ON o.id = cv.organization_id
LEFT JOIN credential_schemas cs ON cs.id = cv.schema_id
LEFT JOIN recipients r ON r.id = cv.recipient_id
WHERE cv.id IN (?)