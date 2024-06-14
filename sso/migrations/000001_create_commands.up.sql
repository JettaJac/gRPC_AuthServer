CREATE TABLE If NOT EXISTS users 
(
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL
    --но лучше так не делать при большом количестве пользователей, лучше вынести в отдельную таблицу
    is_admin BOOLEAN NOT NULL DEFAULT FALSE, 
);

CREATE INDEX IF NOT EXISTSidx_email ON users(email);

CREATE TABLE IF NOT EXISTS apps
(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE,
);

