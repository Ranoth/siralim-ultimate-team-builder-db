# Siralim Ultimate Team Builder Database (Work in Progress)
> [Docker Image](https://hub.docker.com/r/ranoth/siralim-ultimate-team-builder-db)

This repository contains the API for the Siralim Ultimate Team Builder (WIP), a tool designed to help players of the game Siralim Ultimate create and manage their teams effectively. The database includes information on monsters, skills, items, and other relevant data that players can use to optimize their team compositions.

## How to Use
1. Install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/install) if you haven't already.
2. Write a `docker-compose.yaml` file as such:
    ```yaml
    volumes:
        sutbdb_postgres_data:

    services:
        sutbdb:
            image: siralim-ultimate-team-builder-db:latest
            container_name: sutbdb
            environment:
                GOOSE_DBSTRING: host=postgres user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable
            ports:
                - 8080
            depends_on:
            postgres:
                condition: service_healthy
            volumes:
                - ./gameData:/app/gameData:ro
        
        postgres:
            image: postgres:16-alpine
            container_name: sutbdb-postgres
            environment:
                POSTGRES_USER: ${POSTGRES_USER}
                POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
                POSTGRES_DB: ${POSTGRES_DB}
            volumes:
                - sutbdb_postgres_data:/var/lib/postgresql/data
            healthcheck:
                test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
                interval: 5s
                timeout: 5s
                retries: 5
    ```
3. Download game data files from the releases of the [SiralimJSON](https://github.com/iconmaster5326/SiralimJSON) repository, you should download `aggregate.zip` from the latest release and extract the JSON files and `images` folder into the `gameData` folder.
4. Run the the app with `docker compose up -d --build` and the database will be seeded with the latest data from the game. Make sure to set the environment variables in a .env file or in your shell before running the app; the host for the database connection string should be the same as the service name of the postgres container in the docker compose file (in this case, `postgres`).
5. If you want to use the API from the host, you should map the 8080 port of the container to a port on your host machine in the docker compose file.