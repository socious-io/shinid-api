ALTER TABLE credential_verifications
ALTER COLUMN status SET DEFAULT 'CREATED';

UPDATE credential_verifications SET status='CREATED' WHERE status;