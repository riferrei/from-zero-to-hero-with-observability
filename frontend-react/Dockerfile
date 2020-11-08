FROM node:12

ENV ELASTIC_APM_JS_BASE_SERVICE_NAME=frontend-react
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
WORKDIR /usr/src/app
COPY . /usr/src/app

# Metric instrumentation
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/metricbeat/metricbeat-7.9.3-amd64.deb
RUN dpkg -i metricbeat-7.9.3-amd64.deb
RUN cp logs-and-metrics/metricbeat/metricbeat.yml /etc/metricbeat
RUN cp logs-and-metrics/metricbeat/system.yml /etc/metricbeat/modules.d
RUN metricbeat modules enable system

# Network instrumentation
RUN apt update
RUN apt-get install libpcap0.8
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/packetbeat/packetbeat-7.9.3-amd64.deb
RUN dpkg -i packetbeat-7.9.3-amd64.deb
RUN cp logs-and-metrics/packetbeat/packetbeat.yml /etc/packetbeat

# Uptime instrumentation
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/heartbeat/heartbeat-7.9.3-amd64.deb
RUN dpkg -i heartbeat-7.9.3-amd64.deb
RUN cp logs-and-metrics/heartbeat/heartbeat.yml /etc/heartbeat

# Audit instrumentation
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/auditbeat/auditbeat-7.9.3-amd64.deb
RUN dpkg -i auditbeat-7.9.3-amd64.deb
RUN cp logs-and-metrics/auditbeat/auditbeat.yml /etc/auditbeat

RUN npm install --unsafe-perm \
    && npm run-script build
RUN sh install-pkgs.sh

RUN ["chmod", "+x", "/usr/src/app/entrypoint.sh"]
ENTRYPOINT [ "/usr/src/app/entrypoint.sh" ]
