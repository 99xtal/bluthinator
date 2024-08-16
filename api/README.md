# api

This module contains the REST API that the `/web` client uses to access resources from the database, the search engine, and object storage.

## Development Setup

### Backend Setup
Spin up the application using Docker Compose to create the database, search engine, and object storage server.

### Environment
Set your environment variables to point to the servers running in Docker:
```
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=bluthinator
ELASTIC_HOST=localhost
ELASTIC_PORT=9200
ELASTIC_USER=elastic
ELASTIC_PASS=elastic
OBJECT_STORAGE_ENDPOINT=http://localhost:9000
OBJECT_STORAGE_ACCESS_KEY=minio
OBJECT_STORAGE_SECRET_KEY=minio123
```

### Running Locally
Build and run the application
```
go build ./cmd/api/main.go
./main

# or 
go run ./cmd/api/main.go
```
