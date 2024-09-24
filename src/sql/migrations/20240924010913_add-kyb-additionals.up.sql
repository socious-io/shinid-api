CREATE TYPE kyb_verification_status_type AS ENUM ('PENDING', 'APPROVED', 'REJECTED');

ALTER TABLE kyb_verifications
    ADD COLUMN status kyb_verification_status_type DEFAULT 'PENDING' NOT NULL,
    ADD COLUMN created_at TIMESTAMP DEFAULT NOW(),
    ADD COLUMN updated_at TIMESTAMP DEFAULT NOW(),
    DROP COLUMN IF EXISTS description;
