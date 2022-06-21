PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    uuid TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    uuid TEXT PRIMARY KEY,
    content TEXT,
    sender TEXT REFERENCES users,
    recipient TEXT REFERENCES users,
    created_at TEXT NOT NULL
)