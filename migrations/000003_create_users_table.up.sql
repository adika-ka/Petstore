CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    email TEXT NOT NULL,
    password TEXT,
    phone TEXT,
    user_status INT
);