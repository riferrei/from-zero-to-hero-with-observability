#!/bin/bash
# set -x

export BACKEND_URL=http://localhost:8080

sed -i '/REACT_APP_'/d variables.env
echo "REACT_APP_BACKEND_URL=${BACKEND_URL}" >> variables.env

npm run-script start
