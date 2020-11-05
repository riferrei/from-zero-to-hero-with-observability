FROM node:12-alpine

RUN apk update && \
    apk --no-cache add git && \
    rm -rf /var/cache/apk/*

ENV ELASTIC_APM_JS_BASE_SERVICE_NAME=frontend-react
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true

WORKDIR /app
COPY . /app

RUN npm install --unsafe-perm \
    && npm run-script build
RUN sh install-pkgs.sh

CMD ["npm", "run-script", "start"]
