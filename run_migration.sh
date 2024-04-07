#!/bin/bash

DB_URL="postgresql://${MASTER_DB_USER}:${MASTER_DB_PASSWORD}@${MASTER_DB_HOST}:${MASTER_DB_PORT}/${MASTER_DB_NAME}?sslmode=${MASTER_SSL_MODE}"
migrate -path db/migration -database "$DB_URL" -verbose up