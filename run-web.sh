#!/bin/sh


export PORT=8000

export DB_MIGRATION_PATH=file://../migration
export DB_MIGRATION_VERSION=1
export DB_URL=postgres://postgres:pgpassword@localhost:5433/postgres
export DB_USE_SSL=false

go get github.com/githubnemo/CompileDaemon@master

if CompileDaemon --build='go build cmd/web/main.go' --command=./main ; then
    echo "Go build success"
else
    echo "Go build failed"
fi