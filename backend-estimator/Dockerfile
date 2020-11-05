FROM maven:3.6.3-jdk-11

VOLUME /tmp
ADD . /usr/src/app
WORKDIR /usr/src/app
RUN mkdir /etc/app

RUN mvn clean package -DskipTests
RUN cp -v ./target/backend-estimator-service-0.1.0.jar /etc/app/app.jar

ENV AGENT_VERSION=1.18.1
RUN curl -O https://repo1.maven.org/maven2/co/elastic/apm/elastic-apm-agent/$AGENT_VERSION/elastic-apm-agent-$AGENT_VERSION.jar
RUN mv elastic-apm-agent-$AGENT_VERSION.jar /etc/app

ENV JAVA_OPTS="-javaagent:/etc/app/elastic-apm-agent-$AGENT_VERSION.jar"
ENTRYPOINT [ "sh", "-c", "java $JAVA_OPTS -jar /etc/app/app.jar" ]
