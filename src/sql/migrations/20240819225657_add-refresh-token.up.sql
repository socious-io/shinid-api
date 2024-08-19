CREATE TABLE tokens_blacklist (
  id UUID NOT NULL DEFAULT public.uuid_generate_v4() PRIMARY KEY,
  token TEXT UNIQUE,
  expired_at TIMESTAMP DEFAULT NOW()
);