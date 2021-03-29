# build
FROM golang:1.13.8 AS build
ADD . /app
WORKDIR /app
RUN go build -o /out

# deploy
FROM debian:buster
EXPOSE 8080
WORKDIR /
COPY --from=build /out /bin/curfetch
RUN chmod +x /bin/curfetch
CMD ["curfetch", "serve"]