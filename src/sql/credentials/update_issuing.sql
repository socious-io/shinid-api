UPDATE credentials SET
  status=
  issued_at=NOW()
WHERE id=$1
RETURNING *