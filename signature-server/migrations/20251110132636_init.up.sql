CREATE TABLE IF NOT EXISTS licenses (
        id SERIAL PRIMARY KEY,
        user_id TEXT NOT NULL,
        product TEXT NOT NULL,
        issued_at TIMESTAMPTZ NOT NULL,
        expires_at TIMESTAMPTZ NOT NULL,
        hw_fingerprint TEXT NOT NULL,
        nonce TEXT NOT NULL,
        CONSTRAINT uniq_user_fingerprint UNIQUE (user_id, hw_fingerprint)
);
