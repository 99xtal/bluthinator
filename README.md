# Bluthinator 🍌

Bluthinator is an Arrested Development search engine. Search by quotes and generate memes and GIFs from your favorite moments of Seasons 1-3.

![Image of Lucille Bluth, captioned "It's one EC2 instance, Michael. What could it cost, $10?"](https://api.bluthinator.com/caption/S1E06/381291?b=SXQncyBvbmUgRUMyIGluc3RhbmNlLCBNaWNoYWVsLiBXaGF0IGNvdWxkIGl0IGNvc3QsICQxMD8=)

## Project Overview

### Core
`/core` contains the main video processing code which generates search index data and media assets to allow frames to be searchable by their associated quote. [Read More](./core)

### Web Application
`/api` is a RESTful API server which gives access to the following resources defined in `docker-compose.yml`/ [Read More](./api)

- A Postgres database which stores episode metadata
- An ElasticSearch server storing an index of subtitled frames
- An object storage server (Minio) for serving static media assets

`/web` is a Next.js application which serves as the web frontend for Bluthinator. [Read More](./web)

## Development Setup
*Requirements*: Docker, Docker Compose

### Environment
Set the following environment variables:
```
POSTGRES_USER=exampleuser
POSTGRES_PASSWORD=examplepassword
POSTGRES_DB=bluthinator
MINIO_ROOT_USER=minio
MINIO_ROOT_PASSWORD=minio123
ELASTIC_USERNAME=elastic
ELASTIC_PASSWORD=elastic
```

### Running the project
Once you have loaded the episode data to the database, search index, and object storage, you can run only the containers that form the web app:
```
docker compose up -d
```

To shut down running all containers:
```
docker compose down
```
