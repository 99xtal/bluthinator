# Bluthinator

Bluthinator is a search engine for the TV show Arrested Development. It processes a directory of video files containing TV show episodes and extracts unique frame images based on a user-defined threshold. This threshold determines the difference in average color from the previously saved frame, ensuring only significantly different frames are saved. Additionally, it extracts the subtitle for each frame and stores the frame-subtitle pairs in a JSON file.

## Development Setup

### Running the Bluthinator script
Initialize virtual environment
```
python3 -m venv venv
source venv/bin/activate
```

Install dependencies
```
pip install -r requirements.txt
```

Example Usage
```
python3 bluthinator.py ./episodes
python3 bluthinator.py ./episodes -o ./output_dir
```

### Running the entire project
*Requirements*: Docker, Docker Compose

Source the following environment variables:
```
POSTGRES_USER=exampleuser
POSTGRES_PASSWORD=examplepassword
POSTGRES_DB=bluthinator
```

From the root directory, run Docker Compose:
```
docker compose up -d
```

To shut down running containers:
```
docker compose down
```
