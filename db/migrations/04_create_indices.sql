CREATE INDEX idx_frames_episode ON frames(episode);
CREATE INDEX idx_frames_timestamp ON frames(timestamp);

CREATE INDEX idx_episodes_key ON episodes(key);

CREATE INDEX idx_subtitles_episode ON subtitles(episode);

