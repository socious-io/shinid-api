CREATE TYPE credential_status_type AS ENUM ('ISSUED', 'CLAIMED', 'CANCELED');

ALTER TABLE credentials
ADD COLUMN status credential_status_type NOT NULL DEFAULT 'ISSUED',
ADD COLUMN description TEXT,
ADD COLUMN connection_id TEXT,
ADD COLUMN connection_url TEXT,
ADD COLUMN connection_at TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();
