CREATE TABLE frames (
    id SERIAL PRIMARY KEY,
    timestamp INT NOT NULL,
    episode TEXT NOT NULL,
    subtitle TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE episodes (
    id SERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    season INT NOT NULL,
    episode_number INT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ensure the key field in the episodes table is unique
ALTER TABLE episodes
ADD CONSTRAINT unique_key UNIQUE (key);

-- Add a foreign key constraint to the frames table
ALTER TABLE frames
ADD CONSTRAINT fk_episode
FOREIGN KEY (episode) REFERENCES episodes(key);