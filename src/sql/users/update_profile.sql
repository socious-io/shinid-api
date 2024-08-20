UPDATE users
SET 
    first_name=$2,
    last_name=$3,
    bio=$4,
    job_title=$5,
    phone=$6,
    username=$7
WHERE id=$1 RETURNING id