BEGIN;

CREATE TABLE IF NOT EXISTS wallets (
    id      UUID PRIMARY KEY,
    balance DECIMAL NOT NULL CHECK (balance >= 0)
);

COMMIT;
