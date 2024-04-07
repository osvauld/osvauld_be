# Start from golang base image
FROM golang:1.22-alpine as builder

# Install git.
RUN apk update && apk add --no-cache git

# Working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy everythings
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin/migrate

# Copy the Pre-built binary file from the previous stage. Also copy config yml file
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/db/ ./db/

# Copy migration script into the container
COPY --from=builder /app/db/run_migration.sh ./run_migration.sh
RUN chmod +x ./run_migration.sh

# Expose port 8080 to the outside world
EXPOSE 8000

ENTRYPOINT ["sh ./run_migration.sh"]
CMD ["./main"]