#!/bin/bash

heartbeat setup && service heartbeat-elastic start

exec npm run-script start
