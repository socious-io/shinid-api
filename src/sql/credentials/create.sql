INSERT INTO credentials 
  (name, description, schema_id, created_id, organization_id, recipient_id, claims) 
VALUES 
  ($1, $2, $3, $4, $5, $6, $7)
RETURNING *