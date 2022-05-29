# Squzy Storage - open source data storage service

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)

## About

System allows to storage information provided by health checks and agents. It uses clickhouse (optionally postgresSQL).

## API

[**GRPC API**](https://github.com/squzy/squzy_proto/blob/master/proto/v1/squzy_storage.proto#L19) 

## Environment variables

Bold is required

- PORT(9090) - on with port run squzy_storage
- **DB_HOST** - db host
- **DB_PORT** - db port
- **DB_NAME** - db name
- **DB_USER** - db user
- **DB_PASSWORD** - db password
- DB_TYPE - default clickhouse, optionally postgres
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
