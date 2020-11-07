#!/bin/bash

filebeat setup && service filebeat start
metricbeat setup && service metricbeat start
packetbeat setup && service packetbeat start
heartbeat setup && service heartbeat-elastic start
auditbeat setup && service auditbeat start

exec java -javaagent:elastic-apm-agent.jar -jar target/backend-database-service-0.1.0.jar
