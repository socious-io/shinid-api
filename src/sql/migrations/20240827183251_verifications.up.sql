CREATE TABLE credential_verifications (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description TEXT,
  schema_id UUID NOT NULL,
  user_id UUID NOT NULL,
  connection_id TEXT,
  connection_url TEXT,
  body jsonb,
  verified_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_schema FOREIGN KEY (schema_id) REFERENCES credential_schemas(id) ON DELETE CASCADE,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);