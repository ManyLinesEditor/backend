CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE devices (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fingerprint  TEXT        NOT NULL,
    revoked      BOOLEAN     NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, fingerprint)
);

CREATE TABLE features (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE subscriptions (
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feature_id UUID        NOT NULL REFERENCES features(id) ON DELETE CASCADE,
    until      TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (user_id, feature_id)
);

CREATE TABLE payments (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feature_id   UUID        NOT NULL REFERENCES features(id),
    reference    TEXT        NOT NULL UNIQUE,
    status       TEXT        NOT NULL DEFAULT 'pending',
    amount       INTEGER     NOT NULL,
    checkout_url TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE documents (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE deltas (
    update_seq  BIGINT  NOT NULL,
    document_id UUID    NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    device_id   UUID    NOT NULL REFERENCES devices(id),
    payload     BYTEA   NOT NULL,
    PRIMARY KEY (update_seq, document_id)
);

CREATE INDEX idx_devices_user_id ON devices(user_id);
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_deltas_document_id ON deltas(document_id);