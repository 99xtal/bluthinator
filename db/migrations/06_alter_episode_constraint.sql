-- Drop the existing foreign key constraint from the frames table
ALTER TABLE frames
DROP CONSTRAINT fk_episode;

-- Add a new foreign key constraint with ON DELETE NO ACTION
ALTER TABLE frames
ADD CONSTRAINT fk_episode
FOREIGN KEY (episode) REFERENCES episodes(key)
ON DELETE NO ACTION;

-- Drop the existing foreign key constraint from the subtitles table
ALTER TABLE subtitles
DROP CONSTRAINT fk_episode;

-- Add a new foreign key constraint with ON DELETE NO ACTION
ALTER TABLE subtitles
ADD CONSTRAINT fk_episode
FOREIGN KEY (episode) REFERENCES episodes(key)
ON DELETE NO ACTION;