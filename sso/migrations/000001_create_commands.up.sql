-- Active: 1714057451806@@127.0.0.1@5432@sso
CREATE TABLE If NOT EXISTS users (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pass_hash BYTEA NOT NULL,
    --но лучше так не делать при большом количестве пользователей, лучше вынести в отдельную таблицу
    is_admin BOOLEAN NOT NULL DEFAULT FALSE);

CREATE INDEX IF NOT EXISTSidx_email ON users(email);

CREATE TABLE IF NOT EXISTS apps(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE);

