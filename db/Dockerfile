# Use the official PostgreSQL image as a base
FROM postgres:latest

# Set environment variables (these can be overridden by docker-compose)
ENV POSTGRES_USER=${POSTGRES_USER}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
ENV POSTGRES_DB=${POSTGRES_DB}

# Install Python
RUN apt-get update && apt-get install -y python3 python3-pip

# Copy the preprocessing script and the original JSON file into the container
COPY ./scripts /docker-entrypoint-initdb.d/scripts/
COPY ./frame_metadata.json /docker-entrypoint-initdb.d/frame_metadata.json

# Preprocess the frame metadata JSON file for initialization scripts
RUN python3 /docker-entrypoint-initdb.d/scripts/preprocess_json.py /docker-entrypoint-initdb.d/frame_metadata.json

# Copy initialization scripts into the container
COPY ./migrations /docker-entrypoint-initdb.d/