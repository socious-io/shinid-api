UPDATE credentials SET
  name=$2,
  description=$3,
  schema_id=$4,
  claims=$5
WHERE id=$1
RETURNING *