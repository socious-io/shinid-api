UPDATE credential_verifications SET
  present_id=$2,
  updated_at=NOW()
WHERE id=$1
RETURNING *