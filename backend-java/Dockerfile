FROM maven:3.6.3-jdk-11

VOLUME /tmp
ADD . /usr/src/app
WORKDIR /usr/src/app

ENV BEATS_VERSION=7.10.0

# Log instrumentation
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-$BEATS_VERSION-amd64.deb
RUN dpkg -i filebeat-$BEATS_VERSION-amd64.deb
RUN cp logs-and-metrics/filebeat/filebeat.yml /etc/filebeat
RUN cp logs-and-metrics/filebeat/system.yml /etc/filebeat/modules.d
RUN filebeat modules enable system

# Metric instrumentation
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/metricbeat/metricbeat-$BEATS_VERSION-amd64.deb
RUN dpkg -i metricbeat-$BEATS_VERSION-amd64.deb
RUN cp logs-and-metrics/metricbeat/metricbeat.yml /etc/metricbeat
RUN cp logs-and-metrics/metricbeat/system.yml /etc/metricbeat/modules.d
RUN metricbeat modules enable system

# Network instrumentation
RUN apt update
RUN apt-get install libpcap0.8
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/packetbeat/packetbeat-$BEATS_VERSION-amd64.deb
RUN dpkg -i packetbeat-$BEATS_VERSION-amd64.deb
RUN cp logs-and-metrics/packetbeat/packetbeat.yml /etc/packetbeat

# Uptime instrumentation
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/heartbeat/heartbeat-$BEATS_VERSION-amd64.deb
RUN dpkg -i heartbeat-$BEATS_VERSION-amd64.deb
RUN cp logs-and-metrics/heartbeat/heartbeat.yml /etc/heartbeat

# Audit instrumentation
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/auditbeat/auditbeat-$BEATS_VERSION-amd64.deb
RUN dpkg -i auditbeat-$BEATS_VERSION-amd64.deb
RUN cp logs-and-metrics/auditbeat/auditbeat.yml /etc/auditbeat

# App instrumentation
RUN mvn clean package -DskipTests
ENV AGENT_VERSION=1.19.0
RUN curl -O https://repo1.maven.org/maven2/co/elastic/apm/elastic-apm-agent/$AGENT_VERSION/elastic-apm-agent-$AGENT_VERSION.jar
RUN mv elastic-apm-agent-$AGENT_VERSION.jar elastic-apm-agent.jar

RUN ["chmod", "+x", "/usr/src/app/entrypoint.sh"]
ENTRYPOINT [ "/usr/src/app/entrypoint.sh" ]
