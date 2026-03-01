FROM golang:1.25-alpine AS setup
RUN GOBIN=/usr/local/bin go install -tags='no_mysql no_sqlite3 no_ydb' github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

FROM setup AS build
WORKDIR /app

COPY go.mod go.sum ./
COPY . .
RUN go mod download && go build -o /app/sutbdb ./cmd

FROM golang:1.25-alpine AS runtime
WORKDIR /app

ENV GOOSE_DRIVER=postgres
ENV GOOSE_MIGRATION_DIR=/app/migrations

COPY --from=build /app/sutbdb ./
COPY --from=build /app/internal/adapters/postgresql/migrations ./migrations
COPY --from=setup /usr/local/bin/goose /usr/local/bin/goose
COPY entrypoint.sh ./scripts/entrypoint.sh

RUN chmod +x ./scripts/entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["sh", "./scripts/entrypoint.sh"]