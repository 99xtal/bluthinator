#!/bin/bash

# This script is used to extract episode data from TMDB, transform it, and upload it to Postgres.

if [ -z "$TMDB_API_TOKEN" ]; then
    echo "Error: TMDB_API_TOKEN is not set."
    exit 1
fi

TMDB_SEASON_URL="https://api.themoviedb.org/3/tv/4589/season"
SEASONS=(1 2 3)

for SEASON in "${SEASONS[@]}"; do
    curl --request GET \
        --url "$TMDB_SEASON_URL/$SEASON" \
        --header "Authorization: Bearer $TMDB_API_TOKEN" | \
    jq -r '
        def pad_left(s; n; pad): 
            if (n - (s | length)) > 0 then 
                (pad * (n - (s | length))) + s 
            else 
                s 
            end;
        .episodes[] | 
        {
            air_date,
            key: ("S" + (.season_number | tostring) + "E" + (pad_left((.episode_number | tostring); 2; "0"))),
            episode_number, 
            name, 
            overview, 
            season_number,
            directors: (.crew | map(select(.job == "Director") | .name) | join(", ")),
            writers: (.crew | map(select(.job == "Writer") | .name) | join(", "))
        } |
        [.air_date, .key, .episode_number, .name, .overview, .season_number, .directors, .writers] |
        @csv' > episodes.csv

    psql -c "CREATE TABLE tmp_episodes (air_date date, key text, episode_number int, title text, overview text, season int, director text, writer text);"

    psql -c "\COPY tmp_episodes (air_date, key, episode_number, title, overview, season, director, writer) FROM 'episodes.csv' WITH (FORMAT csv, HEADER false);"

    # Update existing rows in the episodes table with the new data
    psql -c "
        UPDATE episodes
        SET air_date = tmp.air_date,
            episode_number = tmp.episode_number,
            title = tmp.title,
            overview = tmp.overview,
            season = tmp.season,
            director = tmp.director,
            writer = tmp.writer
        FROM tmp_episodes tmp
        WHERE episodes.key = tmp.key;
    "

    # Insert new rows into the episodes table
    psql -c "
        INSERT INTO episodes (air_date, key, episode_number, title, overview, season, director, writer)
        SELECT air_date, key, episode_number, title, overview, season, director, writer
        FROM tmp_episodes
        WHERE key NOT IN (SELECT key FROM episodes);
    "

    psql -c "DROP TABLE tmp_episodes;"

    rm episodes.csv
done
