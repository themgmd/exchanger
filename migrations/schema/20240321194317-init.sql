
-- +migrate Up
CREATE TABLE IF NOT EXISTS courses (
    id            serial      primary key,
    currency_from text        not null,
    currency_to   text        not null,
    rate          numeric     not null,
    created_at    timestamptz default now(),
    updated_at    timestamptz default now()
);
-- +migrate Down
