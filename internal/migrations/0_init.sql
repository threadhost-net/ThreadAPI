BEGIN;
CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email         TEXT UNIQUE NOT NULL,
    is_activated  BOOLEAN NOT NULL DEFAULT false,
    is_banned     BOOLEAN NOT NULL DEFAULT false,
    is_admin      BOOLEAN NOT NULL DEFAULT false,
    is_partner    BOOLEAN NOT NULL DEFAULT false,
    balance       BIGINT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS clients (
    client_id     TEXT PRIMARY KEY,
    client_secret TEXT NOT NULL,
    client_name   TEXT NOT NULL,
    redirect_uris TEXT[] NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
COMMIT;