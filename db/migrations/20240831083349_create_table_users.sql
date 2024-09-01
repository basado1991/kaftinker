-- migrate:up
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY,
    username TEXT NOT NULL,
    password BLOB NOT NULL
);

-- migrate:down
DROP TABLE IF EXISTS users;
