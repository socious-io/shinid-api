SELECT 
  cv.*,
  row_to_json(u.*) AS user,
  row_to_json(cs.*) AS schema
FROM credential_verifications cv 
LEFT JOIN users u ON u.id = cv.user_id
LEFT JOIN credential_schemas cs ON cs.id = cv.schema_id
WHERE cv.id IN (?)