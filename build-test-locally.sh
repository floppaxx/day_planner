#!/bin/bash -ex

# This script is used to build the app locally for testing purposes.

docker network create app-test-network || true

docker build -t app:latest .

# Source environment variables from .env file
if [ -f .env.db ]; then
    source .env.db
else
    echo "Error: .env.db file not found!"
    exit 1
fi

docker kill mysql-container || true
docker run -d --rm --name mysql-container --network app-test-network \
    -e MYSQL_ROOT_PASSWORD="$MYSQL_ROOT_PASSWORD" \
    -e MYSQL_DATABASE="$MYSQL_DATABASE" \
    -e MYSQL_USER="$MYSQL_USER" \
    -e MYSQL_PASSWORD="$MYSQL_PASSWORD" \
    -p 3306:3306 mysql:latest

sleep 10

docker run -t --network app-test-network \
    -e DATABASEHOST=mysql-container \
    -e DATABASEPORT=3306 \
    -e DATABASENAME="$MYSQL_DATABASE" \
    -e DATABASEUSER="$MYSQL_USER" \
    -e DATABASEPASSWORD="$MYSQL_PASSWORD" \
    -p 8001:8001 app:latest
