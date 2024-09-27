SELECT id, COUNT(*) OVER () as total_count
	FROM kyb_verifications k
	WHERE k.user_id=$1
LIMIT $2 OFFSET $3