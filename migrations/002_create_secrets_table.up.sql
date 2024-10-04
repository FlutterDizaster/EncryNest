BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS secrets (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(), 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), 
    kind INT NOT NULL,
    data BYTEA NOT NULL,
    CONSTRAINT fk_user
      FOREIGN KEY (user_id) 
      REFERENCES users(id)
      ON DELETE CASCADE
);

COMMIT;
