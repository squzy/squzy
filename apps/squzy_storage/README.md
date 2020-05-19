# Squzy Storage - open source data storage service

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)

## About

System allows to storage information provided by health checks and agents. It uses postgresSQL.

# Usage

## API

[**GRPC API**](https://github.com/squzy/squzy_proto/blob/feat/api-1/proto/v1/squzy_storage.proto#L19) 

# Get requests parameters

The get requests use `pagination` and `time_range` not-required parameters. 
Firstly the `time_range` applied, then `pagination` for all received from database results.

`pagination` is the number of page and number of elements shown on page. If it is nil, all elements will be returned.

For example, the request below will return elements with numbers 60..89.

```shell script
"pagination": {
    "page": 2,
    "limit": 30
},
```

`time_range` is time limit for elements. 

If `from` is nil, elements from the oldest will be taken.

If `to` is nil, elements till the latest will be taken.

```shell script
"time_range": {
    "from": {
      "seconds": 20,
      "nanos": 10
    },
    "to": {
      "seconds": 20,
      "nanos": 10
    }
  }
```

## Environment variables

Bold is required

- PORT(9090) - on with port run squzy
- DB_HOST - postgresSQL host
- DB_PORT - postgresSQL port
- DB_NAME - postgresSQL name
- DB_USER - postgresSQL user
- DB_PASSWORD - postgresSQL password

## Docker

[HUB](https://hub.docker.com/repository/docker/squzy/squzy_monitoring)

For current develop branch use tag: **latest**

Docker Hub

```shell script
docker pull squzy/squzy_storage:v1.6.0
```

### Run locally with docker:

```shell script
docker run -p 9090:9090 squzy/squzy_monitoring:v1.6.0
```

# Want to help?
Want to file a bug, contribute some code, or improve documentation? Excellent!

Add merge request with description.
