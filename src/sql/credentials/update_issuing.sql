UPDATE credentials SET
  status='ISSUED',
  issued_at=NOW()
WHERE id=$1
RETURNING *