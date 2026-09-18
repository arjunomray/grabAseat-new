-- +goose Up
CREATE TABLE event_admins (
    event_id uuid NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (event_id, user_id)
);

CREATE INDEX idx_event_admins_user_id ON event_admins (user_id);

-- Event admins replace the global role.
ALTER TABLE users DROP COLUMN role;

-- +goose Down
ALTER TABLE users ADD COLUMN role text NOT NULL DEFAULT 'user';

DROP INDEX IF EXISTS idx_event_admins_user_id;
DROP TABLE event_admins;
