FROM node:12

ENV ELASTIC_APM_JS_BASE_SERVICE_NAME=frontend-react
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
WORKDIR /usr/src/app
COPY . /usr/src/app

# Uptime instrumentation
RUN curl -L -O https://artifacts.elastic.co/downloads/beats/heartbeat/heartbeat-7.9.3-amd64.deb
RUN dpkg -i heartbeat-7.9.3-amd64.deb
RUN cp logs-and-metrics/heartbeat/heartbeat.yml /etc/heartbeat

RUN npm install --unsafe-perm \
    && npm run-script build
RUN sh install-pkgs.sh

RUN ["chmod", "+x", "/usr/src/app/entrypoint.sh"]
ENTRYPOINT [ "/usr/src/app/entrypoint.sh" ]
