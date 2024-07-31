from airflow.decorators import dag, task
from datetime import datetime
import ffmpeg
import math
from minio import Minio
import os
import numpy as np
from PIL import Image
import psycopg2
from psycopg2 import sql
import pysrt
import re
import utils

@dag(
    schedule_interval=None,
    start_date=datetime(2023, 1, 1),
    catchup=False,
    default_args={'owner': 'airflow', 'retries': 1}
)
def process_episodes():
    @task
    def extract_video_files():
        client = Minio(
            endpoint="minio:9000",
            access_key=os.getenv("MINIO_ACCESS_KEY"),
            secret_key=os.getenv("MINIO_SECRET_KEY"),
            secure=False
        )
        objects = client.list_objects('bluthinator', prefix='episodes/')
        
        local_video_paths = []
        for obj in objects:
            video_path = f'/tmp/{obj.object_name}'
            client.fget_object('bluthinator', obj.object_name, video_path)
            local_video_paths.append(video_path)
        return local_video_paths
    
    @task
    def extract_subtitle_file_from_video(video_file_path):
        subtitle_path = video_file_path.replace('.mkv', '.srt')
        ffmpeg.input(video_file_path).output(subtitle_path, format='srt').run(capture_stdout=True, capture_stderr=True, overwrite_output=True)
        return subtitle_path
    
    @task
    def transform_subtitles(srt_path):
        subs = pysrt.open(srt_path)

        episode_key = os.path.splitext(os.path.basename(srt_path))[0]

        subtitle_records = []
        for sub in subs:
            subtitle_records.append({
                'text': utils.clean_subtitle_text(sub.text),
                'start_timestamp': utils.sub_timestamp_to_ms(sub.start),
                'end_timestamp': utils.sub_timestamp_to_ms(sub.end),
                'episode': episode_key
            })
        
        return subtitle_records
    
    @task
    def load_subtitle_file_to_object_storage(subtitle_file_path):
        print(f'Uploading {subtitle_file_path} to object storage')
        client = Minio(
            endpoint="minio:9000",
            access_key=os.getenv("MINIO_ACCESS_KEY"),
            secret_key=os.getenv("MINIO_SECRET_KEY"),
            secure=False
        )

        bucket_name = 'bluthinator'
        object_name = f'subtitles/{os.path.basename(subtitle_file_path)}'
        
        client.fput_object(
            bucket_name=bucket_name,
            object_name=object_name,
            file_path=subtitle_file_path
        )
    
    @task
    def load_subtitles_to_db(subtitle_records):
        print(f'Loading {len(subtitle_records)} subtitle records to the database')

        # Database connection parameters
        db_params = {
            'dbname': os.getenv('POSTGRES_DB'),
            'user': os.getenv('POSTGRES_USER'),
            'password': os.getenv('POSTGRES_PASSWORD'),
            'host': os.getenv('POSTGRES_HOST', 'db'),
            'port': os.getenv('POSTGRES_PORT', 5432)
        }

        # Connect to the PostgreSQL database
        conn = psycopg2.connect(**db_params)
        cursor = conn.cursor()

        # Delete existing records with the same episode field
        delete_query = sql.SQL("""
            DELETE FROM subtitles WHERE episode = %s
        """)
        cursor.execute(delete_query, (subtitle_records[0]['episode'],))

        # Insert subtitle records into the subtitles table
        insert_query = sql.SQL("""
            INSERT INTO subtitles (episode, text, start_timestamp, end_timestamp) VALUES (%s, %s, %s, %s)
        """)

        for (i, record) in enumerate(subtitle_records):
            print(f'[{i + 1}/{len(subtitle_records)}] {record}')
            cursor.execute(insert_query, (record['episode'], record['text'], record['start_timestamp'], record['end_timestamp']))

        # Commit the transaction and close the connection
        conn.commit()
        cursor.close()
        conn.close()

        return None
    
    @task(multiple_outputs=True)
    def extract_video_frames(video_file_path):
        # Get the dimensions of the video
        frame_width, frame_height, frame_rate = utils.get_video_dimensions(video_file_path)
        frame_size = frame_width * frame_height * 3
        aspect_ratio = frame_width / frame_height;
        frame_number = 0
        frame_metadata = []
        chunk_factor = 10
        threshold = 400
        prev_frame_avg_colors = None
        episode_key = os.path.splitext(os.path.basename(video_file_path))[0]
        output_dir = f'/tmp/frames/{episode_key}'

        # Start the ffmpeg process
        process = ffmpeg.input(video_file_path).output('pipe:', format='rawvideo', pix_fmt='rgb24').run_async(pipe_stdout=True)

        while True:
            in_bytes = process.stdout.read(frame_size)
            if not in_bytes:
                break

            img_array = np.frombuffer(in_bytes, np.uint8).reshape((frame_height, frame_width, 3))
            frame_avg_colors = utils.average_color_per_section(img_array, chunk_factor)
            if prev_frame_avg_colors is None:
                prev_frame_avg_colors = frame_avg_colors
                continue

            diff = utils.color_difference(prev_frame_avg_colors, frame_avg_colors)
            if (diff > threshold):
                timestamp = utils.frame_to_timestamp_ms(frame_number, frame_rate)

                # Convert the raw video frame to a PNG image
                img = Image.frombytes('RGB', (frame_width, frame_height), in_bytes)
                sizes = {
                    "small": 240,
                    "medium": 480,
                    "large": 720,
                }
                for size_name, size in sizes.items():
                    resized_img = img.resize((math.ceil(size * aspect_ratio), size))

                    os.makedirs(f'{output_dir}/{timestamp}', exist_ok=True)
                    resized_img.save(f'{output_dir}/{timestamp}/{size_name}.png')


                frame_metadata.append({
                    'timestamp': timestamp,
                    'episode': episode_key
                })
                prev_frame_avg_colors = frame_avg_colors

            frame_number += 1

        return { 'output_dir': output_dir, 'frame_metadata': frame_metadata }

    @task
    def load_frame_metadata_to_db(frame_metadata: list):
        print(f'Loading {len(frame_metadata)} frame metadata records to the database')

        try:
            # Database connection parameters
            db_params = {
                'dbname': os.getenv('POSTGRES_DB'),
                'user': os.getenv('POSTGRES_USER'),
                'password': os.getenv('POSTGRES_PASSWORD'),
                'host': os.getenv('POSTGRES_HOST', 'db'),
                'port': os.getenv('POSTGRES_PORT', 5432)
            }

            # Connect to the PostgreSQL database
            conn = psycopg2.connect(**db_params)
            cursor = conn.cursor()

            # Delete existing records with the same episode field
            delete_query = sql.SQL("""
                DELETE FROM frames WHERE episode = %s
            """)
            cursor.execute(delete_query, (frame_metadata[0]['episode'],))

            # Insert frame metadata records into the frames table
            insert_query = sql.SQL("""
                INSERT INTO frames (episode, timestamp) VALUES (%s, %s)
            """)

            for (i, record) in enumerate(frame_metadata):
                print(f'[{i + 1}/{len(frame_metadata)}] {record}')
                cursor.execute(insert_query, (record['episode'], record['timestamp']))

            # Commit the transaction and close the connection
            conn.commit()
        except Exception as e:
            print(f'Error loading frame metadata to the database: {e}')
            raise
        finally:
            # Ensure the cursor and connection are closed
            if cursor:
                cursor.close()
            if conn:
                conn.close()

        return None

    @task
    def load_frames_to_object_storage(frame_dir):
        client = Minio(
            endpoint="minio:9000",
            access_key=os.getenv("MINIO_ACCESS_KEY"),
            secret_key=os.getenv("MINIO_SECRET_KEY"),
            secure=False
        )

        bucket_name = 'bluthinator'
        episode_key = os.path.basename(frame_dir)
        objects = []
        for root, dirs, files in os.walk(frame_dir):
            for file in files:
                file_path = os.path.join(root, file)
                object_name = f'frames/{episode_key}/{os.path.relpath(file_path, frame_dir)}'
                objects.append((file_path, object_name))

        for (file_path, object_name) in objects:
            print(f'Uploading {file_path} to object storage')
            client.fput_object(
                bucket_name=bucket_name,
                object_name=object_name,
                file_path=file_path
            )

    video_file_paths = extract_video_files()
    subtitle_file_paths = extract_subtitle_file_from_video.expand(video_file_path=video_file_paths)

    subtitle_records = transform_subtitles.expand(srt_path=subtitle_file_paths)
    load_subtitle_file_to_object_storage.expand(subtitle_file_path=subtitle_file_paths)

    subtitle_records >> load_subtitles_to_db.expand(subtitle_records=subtitle_records)

    outputs = extract_video_frames.expand(video_file_path=video_file_paths)
    frame_metadata_list = outputs.map(lambda x: x['frame_metadata'])
    load_frame_metadata_to_db.expand(frame_metadata=frame_metadata_list)

    frame_dir_list = outputs.map(lambda x: x['output_dir'])
    load_frames_to_object_storage.expand(frame_dir=frame_dir_list)

dag_instance = process_episodes()