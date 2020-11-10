#!/bin/bash
# set -x

export ELASTIC_APM_SERVER_URL=
export ELASTIC_APM_SECRET_TOKEN=
export REDIS_URL=

go build -o backend-golang
./backend-golang >> backend-golang.json
