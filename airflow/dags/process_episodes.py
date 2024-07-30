from airflow.decorators import dag, task
from datetime import datetime

@dag(
    schedule_interval='@daily',
    start_date=datetime(2023, 1, 1),
    catchup=False,
    default_args={'owner': 'airflow', 'retries': 1}
)
def process_episodes():
    @task
    def extract():
        print('extracting')

    @task
    def transform():
        print('transforming')

    @task
    def load():
        print('loading')

    extract() >> transform() >> load()

dag_instance = process_episodes()