# use latest image tag
FROM gradle:latest AS builder

WORKDIR /opt/app

ADD src src
ADD build.gradle settings.gradle .

RUN gradle bootJar --rerun-tasks --no-build-cache

FROM openjdk:17-oracle

# sensitive information in environment variables
ENV DB_USERNAME=postgres
ENV DB_PASSWORD=postgres

COPY --from=builder /opt/app/build/libs/demo-0.0.1-SNAPSHOT.jar app.jar

EXPOSE 8080

# run as root
USER root

ENTRYPOINT ["java", "-jar", "app.jar", "--debug"]