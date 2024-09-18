SELECT * from otps
WHERE user_id=$1 AND expired_at>now()
ORDER BY created_at DESC