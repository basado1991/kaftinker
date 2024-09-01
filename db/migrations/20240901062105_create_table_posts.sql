-- migrate:up
CREATE TABLE posts (
    id INTEGER NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,

    author_id INTEGER NOT NULL,
    FOREIGN KEY(author_id) REFERENCES users(id)
);

-- migrate:down
DROP TABLE IF EXISTS posts;
