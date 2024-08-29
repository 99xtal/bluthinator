#!/bin/bash

# This script is used to process subtitles for episode files and upload the data to Postgres.

# Set PostgreSQL environment variables
# export PGUSER=postgres
# export PGHOST=localhost
# export PGPORT=5432
# export PGPASSWORD=postgres
# export PGDATABASE=bluthinator

# Check if a directory argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <subtitle csv>"
  exit 1
fi

SUBTITLE_CSV_PATH="$1"

# Load the new subtitle records into the database
psql -c "\COPY subtitles (start_timestamp, end_timestamp, text, episode) FROM $SUBTITLE_CSV_PATH WITH (FORMAT csv, HEADER false);"