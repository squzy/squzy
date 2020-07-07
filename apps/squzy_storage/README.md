# Squzy Storage - open source data storage service

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)

## About

System allows to storage information provided by health checks and agents. It uses postgresSQL.

## API

[**GRPC API**](https://github.com/squzy/squzy_proto/blob/master/proto/v1/squzy_storage.proto#L19) 

## Environment variables

Bold is required

- PORT(9090) - on with port run squzy_storage
- **DB_HOST** - postgresSQL host
- **DB_PORT** - postgresSQL port
- **DB_NAME** - postgresSQL name
- **DB_USER** - postgresSQL user
- **DB_PASSWORD** - postgresSQL password
- DB_LOGS(false) - provide logs for DB

## Docker

[HUB](https://hub.docker.com/repository/docker/squzy/squzy_monitoring)

For current develop branch use tag: **latest**

Docker Hub

```shell script
docker pull squzy/squzy_storage:latest
```

# Want to help?
Want to file a bug, contribute some code, or improve documentation? Excellent!

Add merge request with description.
