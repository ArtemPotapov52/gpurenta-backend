CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_id  TEXT NOT NULL UNIQUE,
    email      TEXT NOT NULL,
    name       TEXT NOT NULL,
    picture    TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS agents (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    gpu_model        TEXT NOT NULL,
    vram_gb          INT NOT NULL,
    os               TEXT NOT NULL,
    frp_url          TEXT,
    status           TEXT NOT NULL DEFAULT 'offline',
    supported_images TEXT[] NOT NULL DEFAULT '{}',
    price_per_hour   INT NOT NULL DEFAULT 20,
    last_heartbeat   TIMESTAMPTZ,
    secret           TEXT NOT NULL DEFAULT encode(gen_random_bytes(32), 'hex'),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_agents_status ON agents(status);
CREATE INDEX IF NOT EXISTS idx_agents_owner ON agents(owner_id);

CREATE TABLE IF NOT EXISTS rentals (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id   UUID NOT NULL REFERENCES agents(id),
    renter_id  UUID NOT NULL REFERENCES users(id),
    image      TEXT NOT NULL,
    frp_url    TEXT,
    cost_cents INT NOT NULL DEFAULT 0,
    status     TEXT NOT NULL DEFAULT 'active',
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    ends_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_rentals_renter ON rentals(renter_id);
CREATE INDEX IF NOT EXISTS idx_rentals_agent ON rentals(agent_id);
CREATE INDEX IF NOT EXISTS idx_rentals_status ON rentals(status);

ALTER TABLE rentals ADD COLUMN IF NOT EXISTS access_token TEXT;
CREATE INDEX IF NOT EXISTS idx_rentals_token ON rentals(access_token);
