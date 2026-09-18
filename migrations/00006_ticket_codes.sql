-- +goose Up
ALTER TABLE tickets ADD COLUMN code uuid NOT NULL DEFAULT gen_random_uuid();
CREATE UNIQUE INDEX idx_tickets_code ON tickets (code);

ALTER TABLE tickets ADD COLUMN status text NOT NULL DEFAULT 'valid';
ALTER TABLE tickets ADD COLUMN used_at timestamptz;

-- +goose Down
ALTER TABLE tickets DROP COLUMN used_at;
ALTER TABLE tickets DROP COLUMN status;
DROP INDEX IF EXISTS idx_tickets_code;
ALTER TABLE tickets DROP COLUMN code;
