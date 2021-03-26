# TODO: finish this thing

# Build
FROM golang:1.14-alpine3.11 AS build

WORKDIR /curfetch

COPY . .

RUN apk add build-base
#RUN go build -o ./app .

# Deployment
FROM alpine:3.11
EXPOSE 8080

WORKDIR /app

COPY --from=build /sensor-service/app ./

#ENTRYPOINT [ "./app" ]