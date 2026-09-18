-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL UNIQUE,
    password text NOT NULL
);

CREATE TABLE events (
    id uuid PRIMARY KEY,
    created_by uuid NOT NULL,
    name text NOT NULL,
    time text NOT NULL,
    location text NOT NULL,
    price double precision NOT NULL DEFAULT 0,
    seats_available integer NOT NULL DEFAULT 0,
    tags text[]
);

CREATE TABLE tickets (
    id uuid PRIMARY KEY,
    event_id uuid NOT NULL,
    user_id uuid NOT NULL
);

-- +goose Down
DROP TABLE tickets;
DROP TABLE events;
DROP TABLE users;
