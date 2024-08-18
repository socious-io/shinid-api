CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';

CREATE TABLE media (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  user_id UUID,
  url TEXT NOT NULL UNIQUE,
  filename TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  username VARCHAR(128) UNIQUE NOT NULL,
  password TEXT,
  first_name VARCHAR(128),
  last_name VARCHAR(128),
  email VARCHAR(128) UNIQUE NOT NULL,
  phone VARCHAR(128) UNIQUE,
  job_title VARCHAR(128),
  bio TEXT,
  avatar_id UUID,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_media FOREIGN KEY (avatar_id) REFERENCES media(id) ON DELETE SET NULL
);

CREATE TABLE otps (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  user_id UUID NOT NULL,
  code integer NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  expired_at timestamp with time zone DEFAULT (now() + '00:10:00'::interval) NOT NULL,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE organizations (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  did TEXT,
  name VARCHAR(128),
  description TEXT,
  logo_id UUID,
  is_verified BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_media FOREIGN KEY (logo_id) REFERENCES media(id) ON DELETE SET NULL
);

CREATE TABLE organization_members (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  user_id UUID NOT NULL,
  organization_id UUID NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE (user_id, organization_id),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);

CREATE TABLE kyb_verifications (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  user_id UUID NOT NULL,
  organization_id UUID NOT NULL,
  description TEXT,  
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_org FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);

CREATE TABLE kyb_verification_documents (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  verification_id UUID NOT NULL,
  document UUID NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT fk_verification FOREIGN KEY (verification_id) REFERENCES kyb_verifications(id) ON DELETE CASCADE,
  CONSTRAINT fk_media FOREIGN KEY (document) REFERENCES media(id) ON DELETE CASCADE
);