#!/usr/bin/env bash
FROM --platform=$TARGETPLATFORM eclipse-temurin:21-jre

VOLUME /tmp
COPY maven/${project.build.finalName}.jar /app.jar
RUN sh -x -c 'touch /app.jar'
RUN sh -x -c 'chmod -v +x /app.jar'

ENTRYPOINT ["java", "-Djava.security.egd=file:/dev/./urandom", "-Duser.timezone=GMT+08", "-jar", "/app.jar"]
