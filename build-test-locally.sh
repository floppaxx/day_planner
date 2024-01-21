#!/bin/bash -ex

# This script is used to build the app locally for testing purposes.

docker network create app-test-network || true

docker build -t app:latest .

docker kill mysql-container || true

docker run -d --rm --name mysql-container --network app-test-network \
    -e MYSQL_ROOT_PASSWORD=R00t+ \
    -e MYSQL_DATABASE=app_db \
    -e MYSQL_USER=app_user \
    -e MYSQL_PASSWORD=app_pass123+ \
    -p 3306:3306 mysql:latest

sleep 10

docker run -t --network app-test-network \
    -e DATABASEHOST=mysql-container \
    -e DATABASEPORT=3306 \
    -e DATABASENAME=app_db \
    -e DATABASEUSER=app_user \
    -e DATABASEPASSWORD=app_pass123+ \
    -p 8001:8001 app:latest
