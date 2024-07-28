CREATE TABLE frames (
    id SERIAL PRIMARY KEY,
    frame_number INT NOT NULL,
    subtitle TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);