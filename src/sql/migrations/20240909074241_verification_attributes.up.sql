CREATE TYPE verification_operators_type AS ENUM ('EQUAL', 'NOT', 'BIGGER', 'SMALLER');

CREATE TABLE verification_attribute_values (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  attribute_id UUID NOT NULL,
  schema_id UUID NOT NULL,
  verification_id UUID NOT NULL,
  value TEXT NOT NULL,
  operator verification_operators_type NOT NULL DEFAULT 'EQUAL',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_schema FOREIGN KEY (schema_id) REFERENCES credential_schemas(id) ON DELETE CASCADE,
  CONSTRAINT fk_verification FOREIGN KEY (verification_id) REFERENCES credential_verifications(id) ON DELETE CASCADE,
  CONSTRAINT fk_attribute FOREIGN KEY (attribute_id) REFERENCES credential_attributes(id) ON DELETE CASCADE
);