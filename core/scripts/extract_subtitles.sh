#!/bin/bash

# Check if a directory argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <directory>"
  exit 1
fi

if [ ! -d "$1" ]; then
  echo "Error: $1 is not a valid directory."
  exit 1
fi

EPISODE_DIR="$1"

echo "start_timestamp,end_timestamp,text,episode" > subtitles.csv

for file in "$EPISODE_DIR"/*.mkv; do
  # Extract the episode name without extension
  episode=$(basename "$file" .mkv)
  echo "Processing subtitles for $episode..."

  ffmpeg -i "$file" -map 0:s:0 -f srt - 2>/dev/null | \
  awk -v episode="$(basename "$file" .mkv)" 'BEGIN { RS=""; FS="\n" } 
      {
        # Extract timestamps
        match($2, /([0-9]+):([0-9]+):([0-9]+),([0-9]+) --> ([0-9]+):([0-9]+):([0-9]+),([0-9]+)/, arr);
        start_time_ms = (arr[1] * 3600000) + (arr[2] * 60000) + (arr[3] * 1000) + arr[4];
        end_time_ms = (arr[5] * 3600000) + (arr[6] * 60000) + (arr[7] * 1000) + arr[8];
        
        # Extract subtitle text
        subtitle_text = "";
        for (i=3; i<=NF; i++) {
          subtitle_text = subtitle_text $i " ";
        }

        # Remove newline characters and style tags
        gsub(/\r/, " ", subtitle_text);
        gsub(/\n/, " ", subtitle_text);
        gsub(/<[^>]*>/, "", subtitle_text);
        gsub(/^[ \t]+|[ \t]+$/, "", subtitle_text);  # Trim leading and trailing whitespace
        
        # Escape double quotes by doubling them
        gsub(/"/, "\"\"", subtitle_text);
        
        # Enclose the text in double quotes if it contains a comma or double quote
        if (subtitle_text ~ /[",]/) {
          subtitle_text = "\"" subtitle_text "\"";
        }

        print start_time_ms "," end_time_ms "," subtitle_text "," episode;
      }' >> subtitles.csv
done

