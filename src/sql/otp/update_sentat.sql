UPDATE otps
SET sent_at=(now()+'00:02:00')
WHERE id=$1
RETURNING *