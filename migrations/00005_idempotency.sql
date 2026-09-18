-- +goose Up
CREATE TABLE idempotency_keys (
    key text NOT NULL,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    ticket_id uuid NOT NULL REFERENCES tickets (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, key)
);

-- +goose Down
DROP TABLE idempotency_keys;
