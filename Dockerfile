# Stage 1: Builder
FROM golang:1.22-alpine as builder

# Install git.
RUN apk update && apk add --no-cache git

# Working directory
WORKDIR /app

# Optimize dependency caching:
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o main .

# Stage 2: Migrations
FROM alpine:latest as migrations

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin/migrate

# Copy migration-related files from the builder 
COPY --from=builder /app/db/migration ./db/migration
COPY --from=builder /app/run_migration.sh ./run_migration.sh
RUN chmod +x ./run_migration.sh

# Stage 3: Runtime
FROM alpine:latest as runtime

WORKDIR /root/

# Copy the pre-built binary and templates
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates

# Expose port 8000 to the outside world
EXPOSE 8000

# Entrypoint and Command (from your existing Dockerfile)
ENTRYPOINT ["sh", "./run_migration.sh"] 
CMD ["./main"] 