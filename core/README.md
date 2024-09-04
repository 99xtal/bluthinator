# core
_Requirements_: `ffmpeg`

This module contains tools for processing episode files and generating the metadata and media assets necessary to make Bluthinator work. 

This includes:
- fetching episode metadata from TMDB
- extracting subtitle metadata from the subtitle streams of each episode
- generate image assets and metadata from frames of the video stream that the diffing algorithm deems as visually distinct

![frame from Arrested Development showing "Gob's Program"](https://api.bluthinator.com/caption/S3E11/389875?b=IkdvYidzIFByb2dyYW0iPyBJZGlvdC4)
