CREATE TABLE IF NOT EXISTS licenses (
        id SERIAL PRIMARY KEY,
        fingerprint TEXT NOT NULL UNIQUE,
        product TEXT NOT NULL,
        issued_at TIMESTAMP NOT NULL,
        expires_at TIMESTAMP NOT NULL,
        nonce TEXT NOT NULL,
        activated BOOLEAN NOT NULL DEFAULT FALSE
);
