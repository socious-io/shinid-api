CREATE TABLE credential_schemas (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description TEXT,
  created_id UUID,
  public BOOLEAN DEFAULT false,
  deleteable BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (created_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TYPE attribute_type AS ENUM ('TEXT', 'NUMBER', 'BOOLEAN', 'URL', 'DATETIME', 'EMAIL');

CREATE TABLE credential_attributes (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description TEXT,
  type attribute_type NOT NULL,
  schema_id UUID NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_schema FOREIGN KEY (schema_id) REFERENCES credential_schemas(id) ON DELETE CASCADE
);


CREATE TABLE recipients (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  first_name VARCHAR(128),
  last_name VARCHAR(128),
  email VARCHAR(128) NOT NULL,
  user_id UUID NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE credentials (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,  
  name VARCHAR(128) NOT NULL,
  claims jsonb NOT NULL,
  recipient_id UUID,
  schema_id UUID,
  created_id UUID,
  organization_id UUID,
  expired_at TIMESTAMP,
  issued_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_recipient FOREIGN KEY (recipient_id) REFERENCES recipients(id) ON DELETE SET NULL,
  CONSTRAINT fk_schema FOREIGN KEY (schema_id) REFERENCES credential_schemas(id) ON DELETE SET NULL,
  CONSTRAINT fk_user FOREIGN KEY (created_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);