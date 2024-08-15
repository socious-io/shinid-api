CREATE TYPE user_status AS ENUM ('ACTIVE', 'INACTIVE', 'SUSPENDED');

ALTER TABLE users
ADD COLUMN email_verified_at timestamp DEFAULT NULL,
ADD COLUMN status user_status DEFAULT 'INACTIVE' NOT NULL;

ALTER TABLE otps
ADD COLUMN verified_at timestamp DEFAULT NULL;