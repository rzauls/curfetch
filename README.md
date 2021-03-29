# Curfetch

This is a CLI application that fetches RSS feed containing currency data from a bank
and serves this data as two HTTP endpoints. This was written primarily for educational purposes.

This application is written in go and uses various 3rd party dependencies (modules) for:
 - CLI scaffolding ([spf13/cobra](https://github.com/spf13/cobra)) 
 - parsing rss feeds ([mmcdole/gofeed](https://github.com/mmcdole/gofeed))
 - interacting with cassandra database ([gocql/gocql](https://github.com/gocql/gocql))
 - routing ([gorilla/mux](https://github.com/gorilla/mux))

## Installation
Requirements: 
 - working `docker` and `docker-compose` installation.
 - (for development) working `go` toolchain installation.

1. Clone this repository:
```shell
git clone https://github.com/rzauls/curfetch
cd curfetch
```
2. Build docker images:
```shell
docker-compose build
```
3. Run services (*):
```shell
docker-compose up -d
```
*`cassandra` image needs a bit of time until its reachable from other services, so docker-compose spin up time might take a few seconds
4. Populate database with fresh data:
```shell
docker-compose exec app curfetch fetch
```
5. Access endpoints @ localhost:8080

## Endpoint description

### TODO