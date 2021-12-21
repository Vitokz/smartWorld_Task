CREATE TABLE IF NOT EXISTS users (
    id    SERIAL PRIMARY KEY NOT NULL,
    name varchar NOT NULL DEFAULT '',
    login VARCHAR  NOT NULL UNIQUE DEFAULT '',
    password VARCHAR NOT NULL,
    role  VARCHAR NOT NULL DEFAULT '',
    is_blocked bool NOT NULL DEFAULT false,
    created_at int
)