#!/bin/bash
./cleanup.sh

docker-compose up --build
docker-compose run --rm build-and-copy
./commit-log

./cleanup.sh
