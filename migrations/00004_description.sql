-- +goose Up
ALTER TABLE events ADD COLUMN description text NOT NULL DEFAULT '';

-- The many2many join table was never written to; tickets are the source of
-- truth for registrations.
DROP TABLE IF EXISTS user_events;

-- +goose Down
CREATE TABLE user_events (
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    event_id uuid NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, event_id)
);

ALTER TABLE events DROP COLUMN description;
