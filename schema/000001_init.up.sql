CREATE TABLE IF NOT EXISTS courses (
    id serial not null,
    currency_from text not null,
    currency_to text not null,
    well float,
    updated_at timestamp with time zone
);