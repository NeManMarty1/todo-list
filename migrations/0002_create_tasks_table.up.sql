CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
);