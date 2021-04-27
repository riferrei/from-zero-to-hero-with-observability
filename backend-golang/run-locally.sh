#!/bin/bash
# set -x

export REDIS_URL=localhost:6379
export EXPORTER_ENDPOINT=localhost:8200

go run main.go
