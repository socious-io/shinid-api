UPDATE otps
SET sent_at=now()
WHERE id=$1
RETURNING *