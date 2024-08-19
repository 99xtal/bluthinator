#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <frame_dir>"
  exit 1
fi

frame_dir="$1"
output_csv="output.csv"

# Create or clear the output CSV file
touch "$output_csv"

# Iterate through the directories and extract episode and timestamp
find "$frame_dir" -mindepth 2 -maxdepth 2 -type d | while read -r dir; do
  episode=$(basename "$(dirname "$dir")")
  timestamp=$(basename "$dir")
  echo "$episode,$timestamp" >> "$output_csv"
done

psql -c "\COPY frames (episode, timestamp) FROM '$output_csv' WITH (FORMAT csv, HEADER false);"

rm "$output_csv"

