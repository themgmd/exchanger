CREATE TABLE IF NOT EXISTS courses (
    id serial primary key,
    currency_from text not null,
    currency_to text not null,
    rate numeric not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp
);