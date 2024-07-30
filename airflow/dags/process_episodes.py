from airflow.decorators import dag, task
from datetime import datetime
import ffmpeg
from minio import Minio
import os
import psycopg2
from psycopg2 import sql
import pysrt
import re

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
        ffmpeg.input(video_file_path).output(subtitle_path, format='srt').run(capture_stdout=True, capture_stderr=True)
        return subtitle_path
    
    @task
    def transform_subtitles(srt_path):
        subs = pysrt.open(srt_path)

        episode_key = os.path.splitext(os.path.basename(srt_path))[0]

        subtitle_records = []
        for sub in subs:
            subtitle_records.append({
                'text': clean_subtitle_text(sub.text),
                'start_timestamp': sub_timestamp_to_ms(sub.start),
                'end_timestamp': sub_timestamp_to_ms(sub.end),
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
    
    @task
    def extract_video_frames(video_file_path):
        print(f'Extracting frames from {video_file_path}')
        return None

    video_file_paths = extract_video_files()
    subtitle_file_paths = extract_subtitle_file_from_video.expand(video_file_path=video_file_paths)

    subtitle_records = transform_subtitles.expand(srt_path=subtitle_file_paths)
    load_subtitle_file_to_object_storage.expand(subtitle_file_path=subtitle_file_paths)

    subtitle_records >> load_subtitles_to_db.expand(subtitle_records=subtitle_records)

    video_file_paths >> extract_video_frames.expand(video_file_path=video_file_paths)

dag_instance = process_episodes()