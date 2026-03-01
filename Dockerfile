FROM golang:1.25-alpine AS setup
RUN GOBIN=/usr/local/bin go install github.com/pressly/goose/v3/cmd/goose@latest
WORKDIR /app

FROM setup AS build
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod download && go build -o /app/sutbdb ./cmd

FROM golang:1.25-alpine AS runtime
WORKDIR /app
COPY --from=build /app/sutbdb ./
COPY --from=setup /usr/local/bin/goose /usr/local/bin/goose
COPY --from=build /app/internal/adapters/postgresql/migrations ./internal/adapters/postgresql/migrations
COPY entrypoint.sh ./scripts/entrypoint.sh
RUN chmod +x ./scripts/entrypoint.sh
EXPOSE 8080

ENTRYPOINT ["sh", "./scripts/entrypoint.sh"]