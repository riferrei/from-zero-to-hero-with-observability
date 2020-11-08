#!/bin/bash

metricbeat setup && service metricbeat start
packetbeat setup && service packetbeat start
heartbeat setup && service heartbeat-elastic start
auditbeat setup && service auditbeat start

exec npm run-script start
