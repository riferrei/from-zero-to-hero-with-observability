#!/bin/bash
# set -x

export ELASTIC_APM_SERVER_URL=
export ELASTIC_APM_SECRET_TOKEN=
export ESTIMATOR_URL=

mvn clean package -Dmaven.test.skip=true

AGENT_VERSION=1.18.1
AGENT_FILE=elastic-apm-agent-${AGENT_VERSION}.jar
if [ ! -f "${AGENT_FILE}" ]; then
  curl -O  https://repo1.maven.org/maven2/co/elastic/apm/elastic-apm-agent/${AGENT_VERSION}/elastic-apm-agent-${AGENT_VERSION}.jar
fi

java -javaagent:./${AGENT_FILE} \
-Delastic.apm.server_urls="${ELASTIC_APM_SERVER_URL}" \
-Delastic.apm.secret_token="${ELASTIC_APM_SECRET_TOKEN}" \
-Delastic.apm.service_name="backend-java" \
-Delastic.apm.application_packages="com.riferrei.backend" \
-Delastic.apm.profiling_inferred_spans_enabled=true \
-Delastic.apm.enable_log_correlation=true \
-DESTIMATOR_URL="${ESTIMATOR_URL}" \
-jar target/backend-java-service-0.1.0.jar
