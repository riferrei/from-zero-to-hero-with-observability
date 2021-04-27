#!/bin/bash
# set -x

export ESTIMATOR_URL=http://localhost:8888
mvn clean package -Dmaven.test.skip=true

OTEL_AGENT=opentelemetry-javaagent-all.jar
if [ ! -f "${OTEL_AGENT}" ]; then
  wget -O opentelemetry-javaagent-all.jar https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/download/v1.1.0/opentelemetry-javaagent-all.jar
fi

export EXPORTER_ENDPOINT=http://localhost:8200

java -javaagent:./${OTEL_AGENT} \
-Dotel.traces.exporter=otlp \
-Dotel.metrics.exporter=otlp \
-Dotel.exporter.otlp.endpoint="${EXPORTER_ENDPOINT}" \
-Dotel.resource.attributes=service.name=backend-java,service.version=1.0 \
-DESTIMATOR_URL="${ESTIMATOR_URL}" \
-jar target/backend-java-service-1.0.jar
