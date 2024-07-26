# Bluthinator

Bluthinator is a search engine for the TV show Arrested Development. It processes a directory of video files containing TV show episodes and extracts unique frame images based on a user-defined threshold. The threshold value determines the difference in average color from the previously saved frame, ensuring only significantly different frames are saved.

## Development Setup
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
```