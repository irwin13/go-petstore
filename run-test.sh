#!/bin/sh

# if you need to build / rebuild docker image
#docker-compose up --build &
#sleep 5

docker-compose up &

# waiting for db
sleep 3 

# run all test
export DB_MIGRATION_PATH=file://../../migration
export DB_MIGRATION_VERSION=1
export DB_URL=postgres://postgres:pgpassword@localhost:5433/postgres
export DB_USE_SSL=false

go test -v ./...

# run test per package
#go test -v ./internal/dao/...
#go test -v ./internal/service/...

docker-compose stop