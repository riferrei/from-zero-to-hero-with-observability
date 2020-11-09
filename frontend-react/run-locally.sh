#!/bin/bash
# set -x

export ELASTIC_APM_SERVER_URL=
export BACKEND_URL=

sed -i '/REACT_APP_'/d variables.env
echo "REACT_APP_BACKEND_URL=${BACKEND_URL}" >> variables.env
echo "REACT_APP_APM_SERVER_URL=${ELASTIC_APM_SERVER_URL}" >> variables.env

npm run-script start
