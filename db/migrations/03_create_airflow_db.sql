CREATE DATABASE airflow;

-- Create user 'airflow' with a password
CREATE USER airflow WITH PASSWORD 'airflow';

-- Grant all privileges on the 'airflow' database to the 'airflow' user
GRANT ALL PRIVILEGES ON DATABASE airflow TO airflow;

-- Grant privileges on the public schema of the 'airflow' database to the 'airflow' user
\connect airflow
GRANT ALL PRIVILEGES ON SCHEMA public TO airflow;