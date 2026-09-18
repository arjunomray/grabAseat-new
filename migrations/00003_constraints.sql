-- +goose Up
ALTER TABLE events
    ADD CONSTRAINT events_created_by_fkey
        FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE tickets
    ADD CONSTRAINT tickets_event_id_fkey
        FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE,
    ADD CONSTRAINT tickets_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

CREATE INDEX idx_events_created_by ON events (created_by);
CREATE INDEX idx_events_time ON events (time);
CREATE INDEX idx_tickets_event_id ON tickets (event_id);
CREATE INDEX idx_tickets_user_id ON tickets (user_id);

-- A user can hold at most one ticket per event.
CREATE UNIQUE INDEX idx_tickets_event_id_user_id ON tickets (event_id, user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_tickets_event_id_user_id;
DROP INDEX IF EXISTS idx_tickets_user_id;
DROP INDEX IF EXISTS idx_tickets_event_id;
DROP INDEX IF EXISTS idx_events_time;
DROP INDEX IF EXISTS idx_events_created_by;

ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_user_id_fkey;
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_event_id_fkey;
ALTER TABLE events DROP CONSTRAINT IF EXISTS events_created_by_fkey;
