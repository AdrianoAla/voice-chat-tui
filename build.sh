#!/bin/bash

# build.sh

if [ "$1" == "server" ]; then
    go build -o bin/server ./server/main.go 
    if [ "$2" == "run" ]; then
        ./bin/server
    fi
elif [ "$1" == "client" ]; then
    go build -o bin/client ./client/main.go 
    if [ "$2" == "run" ]; then
        ./bin/client
    fi
else
    echo "Usage: $0 [client|server]"
    exit 1
fi
