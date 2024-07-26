import chardet
import ffmpeg
import json
import math
import pysrt
import os
import numpy as np
from PIL import Image

def get_video_dimensions(video_path):
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

def frame_to_timestamp(frame_number, frame_rate):
    total_seconds = frame_number / frame_rate
    hours = int(total_seconds // 3600)
    minutes = int((total_seconds % 3600) // 60)
    seconds = int(total_seconds % 60)
    milliseconds = int((total_seconds - int(total_seconds)) * 1000)
    return f"{hours:02}:{minutes:02}:{seconds:02},{milliseconds:03}"

def detect_encoding(file_path):
    with open(file_path, 'rb') as f:
        raw_data = f.read(10000)  # Read the first 10,000 bytes
    result = chardet.detect(raw_data)
    return result['encoding']

def get_subtitle_for_frame(srt_path, frame_number, frame_rate):
    encoding = detect_encoding(srt_path)
    subs = pysrt.open(srt_path, encoding=encoding)
    timestamp = frame_to_timestamp(frame_number, frame_rate)
    for sub in subs:
        if sub.start <= timestamp <= sub.end:
            return sub.text
    return None

def extract_and_save_frames(video_path, srt_path: str, output_dir, threshold=500, chunk_factor=10):
    # Get the dimensions of the video
    frame_width, frame_height, frame_rate = get_video_dimensions(video_path)
    frame_size = frame_width * frame_height * 3
    frame_number = 0
    prev_frame_avg_colors = None

    frame_metadata = []

    # Start the ffmpeg process
    process = ffmpeg.input(video_path).output('pipe:', format='rawvideo', pix_fmt='rgb24').run_async(pipe_stdout=True)

    while True:
        in_bytes = process.stdout.read(frame_size)
        if not in_bytes:
            break

        img_array = np.frombuffer(in_bytes, np.uint8).reshape((frame_height, frame_width, 3))
        frame_avg_colors = average_color_per_section(img_array, chunk_factor)
        if prev_frame_avg_colors is None:
            prev_frame_avg_colors = frame_avg_colors
            continue

        diff = color_difference(prev_frame_avg_colors, frame_avg_colors)
        if (diff > threshold):
            # Convert the raw video frame to a PNG image
            img = Image.frombytes('RGB', (frame_width, frame_height), in_bytes)
            img.save(f'{output_dir}/frames/{frame_number:04d}.png')

            # Get the subtitle for the frame
            subtitle = get_subtitle_for_frame(srt_path, frame_number, frame_rate)
            print(subtitle);

            frame_metadata.append({
                'frame_number': frame_number,
                'subtitle': subtitle,
            })
            prev_frame_avg_colors = frame_avg_colors

        frame_number += 1

    # Write the list to a JSON file
    with open(f'{output_dir}/frame_metadata.json', 'w') as json_file:
        json.dump(frame_metadata, json_file, indent=4)

    process.wait()

def main():
    input_video = './episodes/S1E01.mkv'
    input_srt = './episodes/S1E01.srt'
    output_dir = './output/S1E01'

    # Create the output directory if it does not exist
    os.makedirs(f'{output_dir}/frames', exist_ok=True)

    extract_and_save_frames(input_video, input_srt, output_dir)

if __name__ == "__main__":
    main()