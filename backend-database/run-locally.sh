#!/bin/bash
# set -x

AGENT_VERSION=1.18.1
AGENT_FILE=elastic-apm-agent-${AGENT_VERSION}.jar
APM_SERVER_URL=
APM_SECRET_TOKEN=

if [ ! -f "${AGENT_FILE}" ]; then
  curl -O  https://repo1.maven.org/maven2/co/elastic/apm/elastic-apm-agent/${AGENT_VERSION}/elastic-apm-agent-${AGENT_VERSION}.jar
fi

mvn clean package -Dmaven.test.skip=true

java -javaagent:./${AGENT_FILE} \
-Delastic.apm.server_urls="${APM_SERVER_URL}" \
-Delastic.apm.secret_token="${APM_SECRET_TOKEN}" \
-Delastic.apm.service_name="backend-database" \
-Delastic.apm.application_packages="com.riferrei.database" \
-Delastic.apm.profiling_inferred_spans_enabled=true \
-Delastic.apm.enable_log_correlation=true \
-DBACKEND_ESTIMATOR_URI=http://localhost:8888/estimateValue \
-jar target/backend-database-service-0.1.0.jar
