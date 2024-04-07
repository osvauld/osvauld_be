#!/bin/bash

db_url="postgresql://${MASTER_DB_USER}:${MASTER_DB_PASSWORD}@${MASTER_DB_HOST}:${MASTER_DB_PORT}/${MASTER_DB_NAME}?sslmode=${MASTER_SSL_MODE}"
make migrateup DB_URL="$db_url"