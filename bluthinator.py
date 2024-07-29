import argparse
import chardet
import ffmpeg
import glob
import json
import math
import numpy as np
import os
from PIL import Image
import pysrt
import re

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

def frame_to_timestamp_ms(frame_number, frame_rate):
    return frame_number * 1000 // frame_rate

def ms_to_sub_timestamp(ms):
    hours = ms // 3600000
    minutes = (ms % 3600000) // 60000
    seconds = (ms % 60000) // 1000
    milliseconds = ms % 1000
    return f"{hours:02}:{minutes:02}:{seconds:02},{milliseconds:03}"

def detect_encoding(file_path):
    with open(file_path, 'rb') as f:
        raw_data = f.read(10000)  # Read the first 10,000 bytes
    result = chardet.detect(raw_data)
    return result['encoding']

def parse_subtitles(file_path) -> pysrt.SubRipFile:
    encoding = detect_encoding(file_path)
    return pysrt.open(file_path, encoding=encoding)

def find_subtitle_for_timestamp(subs: pysrt.SubRipFile, timestamp: int) -> str:
    sub_timestamp = ms_to_sub_timestamp(timestamp)
    for sub in subs:
        if sub.start <= sub_timestamp <= sub.end:
            return sub.text
    return None

def clean_subtitle_text(subtitle):
    # Remove style tags (e.g., <i>...</i>)
    subtitle = re.sub(r'<[^>]+>', '', subtitle)
    # Replace newline characters with spaces
    subtitle = subtitle.replace('\n', ' ')
    # Remove extra spaces
    subtitle = ' '.join(subtitle.split())
    return subtitle

def extract_and_save_frames(video_path, output_dir, threshold=400, chunk_factor=10):
    # Get the dimensions of the video
    frame_width, frame_height, frame_rate = get_video_dimensions(video_path)
    frame_size = frame_width * frame_height * 3
    aspect_ratio = frame_width / frame_height;
    frame_number = 0
    prev_frame_avg_colors = None

    frame_metadata = []

    # Extract the subtitles from the video
    subtitle_path = f'{output_dir}/subtitles.srt'
    ffmpeg.input(video_path).output(subtitle_path, format='srt').run(overwrite_output=True)
    subs = parse_subtitles(subtitle_path)

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
            timestamp = frame_to_timestamp_ms(frame_number, frame_rate)

            # Convert the raw video frame to a PNG image
            img = Image.frombytes('RGB', (frame_width, frame_height), in_bytes)
            sizes = {
                "small": 240,
                "medium": 480,
                "large": 720,
            }
            for size_name, size in sizes.items():
                resized_img = img.resize((math.ceil(size * aspect_ratio), size))

                os.makedirs(f'{output_dir}/frames/{timestamp}', exist_ok=True)
                resized_img.save(f'{output_dir}/frames/{timestamp}/{size_name}.png')

            # Get the subtitle for the frame
            subtitle = find_subtitle_for_timestamp(subs, timestamp)

            frame_metadata.append({
                'timestamp': timestamp,
                'subtitle': clean_subtitle_text(subtitle) if subtitle else None,
            })
            prev_frame_avg_colors = frame_avg_colors

        frame_number += 1

    process.kill()

    return frame_metadata

def main():
    parser = argparse.ArgumentParser(description='Extract frames and metadata from video files')
    parser.add_argument('input_dir', type=str, help='Input directory containing video files')
    parser.add_argument('-o', '--output', type=str, default='./output', help='Base directory for output (default: ./output)')

    # Parse the command-line arguments
    args = parser.parse_args()
    input_dir = args.input_dir
    output_base_dir = args.output

    # Iterate over all video files in the input directory
    video_files = glob.glob(os.path.join(input_dir, '*.mkv'))

    combined_frame_metadata = []

    for video_file in video_files:
        # Extract the base name of the video file (e.g., S1E01 from S1E01.mkv)
        base_name = os.path.splitext(os.path.basename(video_file))[0]
        output_dir = os.path.join(output_base_dir, base_name)

        # Create the output directory for the current video file
        os.makedirs(f'{output_dir}/frames', exist_ok=True)

        # Run the extract function on the current video file
        frame_metadata = extract_and_save_frames(video_file, output_dir)

        # Add the episode name to the metadata
        for entry in frame_metadata:
            entry['episode'] = base_name
        combined_frame_metadata.extend(frame_metadata)

    # Write the list to a JSON file
    with open(f'{output_base_dir}/frame_metadata.json', 'w') as json_file:
        json.dump(combined_frame_metadata, json_file, indent=4)

if __name__ == "__main__":
    main()