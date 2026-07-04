CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

CREATE TABLE orders (
    number TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    status TEXT NOT NULL,
    accrual DOUBLE PRECISION,
    uploaded_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE withdrawals (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    order_number TEXT NOT NULL,
    sum DOUBLE PRECISION NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_orders_user_id_uploaded_at
    ON orders(user_id, uploaded_at DESC);

CREATE INDEX idx_withdrawals_user_id_processed_at
    ON withdrawals(user_id, processed_at DESC);
