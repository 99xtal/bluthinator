# Bluthinator

Bluthinator is a search engine for the TV show Arrested Development. It processes a directory of video files containing TV show episodes and extracts unique frame images based on a user-defined threshold. This threshold determines the difference in average color from the previously saved frame, ensuring only significantly different frames are saved. Additionally, it extracts the subtitle for each frame and stores the frame-subtitle pairs in a JSON file.

## Development Setup

### Running the entire project
*Requirements*: Docker, Docker Compose

Source the following environment variables:
```
POSTGRES_USER=exampleuser
POSTGRES_PASSWORD=examplepassword
POSTGRES_DB=bluthinator
MINIO_ROOT_USER=minio
MINIO_ROOT_PASSWORD=minio123
```

From the root directory, run Docker Compose:
```
docker compose up -d
```

To shut down running containers:
```
docker compose down
```
