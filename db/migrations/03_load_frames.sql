-- Create a temporary table to hold the JSON frame data
CREATE TABLE temp_json_data (
    data JSONB
);

-- Load the JSON data into the temporary table
COPY temp_json_data(data)
FROM '/docker-entrypoint-initdb.d/frame_metadata.json';

-- Insert data into the frames table
INSERT INTO frames (timestamp, episode, subtitle)
SELECT
    (data->>'timestamp')::INT AS timestamp,
    data->>'episode',
    data->>'subtitle'
FROM temp_json_data;

-- Drop the temporary table
DROP TABLE temp_json_data;