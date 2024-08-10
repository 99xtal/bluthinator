import ffmpeg
import math
import numpy as np

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
