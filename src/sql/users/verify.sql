UPDATE users
SET status=$2, email_verified_at=NOW()
WHERE id=$1