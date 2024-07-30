import ffmpeg
import math
import pysrt
import numpy as np
import re

def get_video_dimensions(video_path: str):
    # Use ffprobe to get the metadata of the video file
    probe = ffmpeg.probe(video_path)
    # Extract the video stream information
    video_stream = next((stream for stream in probe['streams'] if stream['codec_type'] == 'video'), None)
    if video_stream is None:
        raise ValueError('No video stream found')
    width = int(video_stream['width'])
    height = int(video_stream['height'])
    raw_frame_rate = video_stream['avg_frame_rate']
    numerator, denominator = map(int, raw_frame_rate.split('/'))
    frame_rate = math.ceil(numerator / denominator)
    return width, height, frame_rate

def average_color_per_section(img_array: np.array, sections: int) -> list:
    width, height, _ = img_array.shape
    section_width = width // sections
    section_height = height // sections

    avg_colors = []

    for i in range(sections):
        for j in range(sections):
            section = img_array[i * section_width: (i + 1) * section_width, j * section_height: (j + 1) * section_height]
            avg_color = np.mean(section, axis=(0, 1))
            avg_colors.append(avg_color)

    return np.array(avg_colors)

def color_difference(avg_colors_1, avg_colors_2):
    return np.linalg.norm(avg_colors_1 - avg_colors_2)

def frame_to_timestamp_ms(frame_number, frame_rate):
    return frame_number * 1000 // frame_rate

def sub_timestamp_to_ms(timestamp: pysrt.SubRipTime) -> int:
    # Split the timestamp into hours, minutes, seconds, and milliseconds
    timestamp_str = str(timestamp)
    hours, minutes, seconds_milliseconds = timestamp_str.split(':')
    seconds, milliseconds = seconds_milliseconds.split(',')
    
    # Convert each part to milliseconds
    hours_ms = int(hours) * 3600000
    minutes_ms = int(minutes) * 60000
    seconds_ms = int(seconds) * 1000
    milliseconds = int(milliseconds)
    
    # Sum all parts to get the total milliseconds
    total_ms = hours_ms + minutes_ms + seconds_ms + milliseconds
    return total_ms

def clean_subtitle_text(subtitle):
    # Remove style tags (e.g., <i>...</i>)
    subtitle = re.sub(r'<[^>]+>', '', subtitle)
    # Replace newline characters with spaces
    subtitle = subtitle.replace('\n', ' ')
    # Remove extra spaces
    subtitle = ' '.join(subtitle.split())
    return subtitle