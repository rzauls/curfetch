# Curfetch

This is a CLI application that fetches RSS feed containing currency data and serves this data as two HTTP endpoints. This was written primarily for educational purposes.

This application is written in go and uses various 3rd party dependencies (modules) for:
 - CLI scaffolding ([spf13/cobra](https://github.com/spf13/cobra)) 
 - parsing rss feeds ([mmcdole/gofeed](https://github.com/mmcdole/gofeed))
 - interacting with cassandra database ([gocql/gocql](https://github.com/gocql/gocql))
 - routing ([gorilla/mux](https://github.com/gorilla/mux))

Since this is an example application, not meant for production, the default cassandra DB keyspace has only 1 node, a replication factor of 1 and uses the default credentials.

`curfetch` has 2 main CLI commands:
    
- `fetch` - fetch and parse RSS feed, insert parsed data in DB (ideally used with `cron` or other task scheduler)
- `serve` - serve data over http

To see detailed instructions about each commands use `--help` flag

To execute any command inside docker use:
```shell
docker-compose exec app curfetch <command> <flags>
```

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
*`cassandra` image needs a bit of time until its reachable from other services, so docker-compose spin-up time might take a few seconds

4. Populate database with fresh data:
```shell
docker-compose exec app curfetch fetch
```
5. Access endpoints @ localhost:8080


If you wish to use `curfetch` without docker, or within another docker-compose setup, use following environment variables:
```
  CASS_HOST: cassandra
  CASS_KEYSPACE: curfetch // modifying keyspace would require modifying the cassandra_schema.cql script aswell
  CASS_USERNAME: cassandra
  CASS_PASSWORD: cassandra
```


## Endpoint description

### Show all latest currency values

**URL** : `/newest`

**Method** : `GET`

### Success Response

**Code** : `200 OK`

**Example**

```json
[
  {
      "code":"CNY",
      "value":"7.73340000",
      "pub_date":"2021-03-29T00:00:00Z"
  },
  {
      "code":"AUD",
      "value":"1.53980000",
      "pub_date":"2021-03-29T00:00:00Z"
  }
]
```

### Show all data points for specific currency

**URL** : `/history/{currency_code}`

**Method** : `GET`

### Success Response

**Code** : `200 OK`

**Example `/history/usd`**

```json
[
  {
    "code":"USD",
    "value":"1.18250000",
    "pub_date":"2021-03-24T00:00:00Z"
  },
  {
    "code":"USD",
    "value":"1.18020000",
    "pub_date":"2021-03-25T00:00:00Z"
  },
  {
    "code":"USD",
    "value":"1.17820000",
    "pub_date":"2021-03-26T00:00:00Z"
  },
  {
    "code":"USD",
    "value":"1.17840000",
    "pub_date":"2021-03-29T00:00:00Z"
  }
]
```
