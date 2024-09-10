#FROM gradle:4.7.0-jdk8-alpine AS build
#COPY --chown=gradle:gradle . /home/gradle/src
#WORKDIR /home/gradle/src
#RUN gradle build --no-daemon
#
#FROM openjdk:8-jre-slim
#
#EXPOSE 8080
#
#RUN mkdir /app
#
#COPY --from=build /home/gradle/src/build/libs/*.jar /app/spring-boot-application.jar
#
#ENTRYPOINT ["java", "-XX:+UnlockExperimentalVMOptions", "-XX:+UseCGroupMemoryLimitForHeap", "-Djava.security.egd=file:/dev/./urandom","-jar","/app/spring-boot-application.jar"]

FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./

RUN #go mod download

COPY . .

RUN go build -o main ./cmd/main.go

#FROM alpine:latest
#
#WORKDIR /root/
#
#COPY --from=build /app/main .
#
EXPOSE 8080

CMD ["./main"]
