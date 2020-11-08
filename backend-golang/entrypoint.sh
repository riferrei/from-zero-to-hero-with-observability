#!/bin/bash

filebeat setup && service filebeat start
metricbeat setup && service metricbeat start
packetbeat setup && service packetbeat start
heartbeat setup && service heartbeat-elastic start
auditbeat setup && service auditbeat start

exec /usr/src/app/backend-golang >> /usr/src/app/backend-golang.json
