# Siralim Ultimate Team Builder Database

This repository contains the database for the Siralim Ultimate Team Builder, a tool designed to help players of the game Siralim Ultimate create and manage their teams effectively. The database includes information on monsters, skills, items, and other relevant data that players can use to optimize their team compositions.

## How to Use
1. Clone the repository to your local machine.
2. Install goose and sqlc.
3. write a docker-compose.yaml file to run the database and dbseeder.
    ```yaml
    volumes:
        sutbdb_postgres_data:

    services:
        sutbdb:
            build: .
            container_name: sutbdb
            environment:
            GOOSE_DBSTRING: host=postgres user=sutbdb password=sutbdb dbname=sutbdb sslmode=disable
            ports:
            - "8080:8080"
            depends_on:
            postgres:
                condition: service_healthy
        
        postgres:
            image: postgres:16-alpine
            container_name: sutbdb-postgres
            environment:
            POSTGRES_USER: sutbdb
            POSTGRES_PASSWORD: sutbdb
            POSTGRES_DB: sutbdb
            ports:
            - "5432:5432"
            volumes:
            - sutbdb_postgres_data:/var/lib/postgresql/data
            healthcheck:
            test: ["CMD-SHELL", "pg_isready -U sutbdb -d sutbdb"]
            interval: 5s
            timeout: 5s
            retries: 5
    ```