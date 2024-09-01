CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY,
    username TEXT NOT NULL,
    password BLOB NOT NULL
);
CREATE TABLE posts (
    id INTEGER NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,

    author_id INTEGER NOT NULL,
    FOREIGN KEY(author_id) REFERENCES users(id)
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20240831083349'),
  ('20240901062105');
