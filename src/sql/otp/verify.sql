UPDATE otps
SET verified_at=now()
WHERE user_id=$1 AND code=$2 AND verified_at=NULL