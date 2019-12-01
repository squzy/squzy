# Squzy - opensource monitoring system
[![codecov](https://codecov.io/gh/squzy/squzy/branch/develop/graph/badge.svg)](https://codecov.io/gh/squzy/squzy)
[![GolangCI](https://golangci.com/badges/github.com/squzy/golangci-lint.svg)](https://golangci.com)

## About

Squzy - is a high-performance open-source monitoring system that supports

**System Health Checks Capabilities**
1) HTTP/HTTPS
2) TCP
3) GRPC - https://github.com/grpc/grpc/blob/master/doc/health-checking.md
4) SiteMap.xml - https://www.sitemaps.org/protocol.html

# Usage

## API
Squzy server implement GRPC API. 

https://github.com/squzy/squzy_proto/blob/master/proto/v1/server.proto

## Storage
By default squzy use stdout for logs, but can be configured by ENV.

Storage should implement that API:

https://github.com/squzy/squzy_proto/blob/master/proto/v1/storage.proto

###Environment variables
- PORT(8080) - on with port run squzy
- STORAGE_HOST - log storage host(example *localhost:9090*)
- STORAGE_TIMEOUT - timeout for connect to log storage

## Docker

```shell script
docker pull docker.pkg.github.com/squzy/squzy/squzy_app:develop
```

# Want to help?
Want to file a bug, contribute some code, or improve documentation? Excellent!

Add merge request with description.