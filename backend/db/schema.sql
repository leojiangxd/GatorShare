DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS comments;

CREATE TABLE users (
    user_id INTEGER PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    salt TEXT NOT NULL,
    hashed_password TEXT NOT NULL
);

CREATE TABLE posts (
    post_id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES useres(user_id),
    title TEXT,
    content TEXT,
    created_at TIMESTAMP,
    likes INTEGER,
    dislikes INTEGER
);

CREATE TABLE comments (
    comment_id INTEGER PRIMARY KEY,
    post_id INTEGER REFERENCES posts(post_id),
    user_id INTEGER REFERENCES users(user_id),
    content TEXT,
    created_at TIMESTAMP,
    likes INTEGER,
    dislikes INTEGER
);