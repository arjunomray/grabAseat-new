-- +goose Up
ALTER TABLE events
    ALTER COLUMN time TYPE timestamptz USING time::timestamptz,
    ALTER COLUMN price TYPE bigint USING (price * 100)::bigint;

ALTER TABLE events RENAME COLUMN price TO price_cents;

ALTER TABLE users
    ADD COLUMN created_at timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN updated_at timestamptz NOT NULL DEFAULT now();

ALTER TABLE events
    ADD COLUMN created_at timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN updated_at timestamptz NOT NULL DEFAULT now();

ALTER TABLE tickets
    ADD COLUMN created_at timestamptz NOT NULL DEFAULT now(),
    ADD COLUMN updated_at timestamptz NOT NULL DEFAULT now();

-- +goose Down
ALTER TABLE tickets
    DROP COLUMN updated_at,
    DROP COLUMN created_at;

ALTER TABLE events
    DROP COLUMN updated_at,
    DROP COLUMN created_at;

ALTER TABLE users
    DROP COLUMN updated_at,
    DROP COLUMN created_at;

ALTER TABLE events RENAME COLUMN price_cents TO price;
ALTER TABLE events
    ALTER COLUMN price TYPE double precision USING (price / 100.0),
    ALTER COLUMN time TYPE text USING time::text;
