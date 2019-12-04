# Squzy - opensource monitoring system
[![codecov](https://codecov.io/gh/squzy/squzy/branch/develop/graph/badge.svg)](https://codecov.io/gh/squzy/squzy)
[![GolangCI](https://golangci.com/badges/github.com/squzy/golangci-lint.svg)](https://golangci.com)

## About

Squzy - is a high-performance open-source monitoring system written in Golang with [Bazel](https://bazel.build/) and love.

**System Health Checks Capabilities**
1) HTTP/HTTPS
2) TCP
3) GRPC - https://github.com/grpc/grpc/blob/master/doc/health-checking.md
4) SiteMap.xml - https://www.sitemaps.org/protocol.html

# Usage

## API
Squzy server implement [GRPC API](https://github.com/squzy/squzy_proto/blob/master/proto/v1/server.proto). 

https://github.com/squzy/squzy_proto/blob/master/proto/v1/server.proto

## Examples of call from [BloomRPC](https://github.com/uw-labs/bloomrpc)

### Http/Https check:

Usually that check used for monitoring web sites

```shell script
{
  "interval": 10, - 10 second interval
  "http_check": {
    "method": "GET", - method GET/POST/PUT/DELETE/HEAD
    "url": "https://google.com", - url which should call
    "headers": {
      "custom": "yes",
    },
    "statusCode": 200 - expected statusCode
  }
}
```

### Tcp check:

Check good use for monitoring open ports or not

```shell script
{
  "interval": 10, - 10 second interval
  "tcp_check": {
    "host": "localhost", - host
    "port": 6345 - port
  },
}
```

### SiteMap check:

That check good usage when you have critical URL in sitemap, if any of URL throw error check will be failed

```shell script
{
  "interval": 10,
  "sitemap_check": {
    "url": "https://www.sitemaps.org/sitemap.xml" - url of sitemap (https://www.sitemaps.org/sitemap.xml)
  },
}
```

### GRPC check:

Check better to use for internal testing of API services

```shell script
{
  "interval": 10,
  "grpc_check": {
    "service": "Check", - service name
    "host": "localhost", - host
    "port": 9090 - port
  },
}
```

## Storage
By default squzy use stdout for logs, but can be configured by ENV.

Storage should implement that [API](https://github.com/squzy/squzy_proto/blob/master/proto/v1/storage.proto):

https://github.com/squzy/squzy_proto/blob/master/proto/v1/storage.proto

## Environment variables
- PORT(8080) - on with port run squzy
- STORAGE_HOST - log storage host(example *localhost:9090*)
- STORAGE_TIMEOUT - timeout for connect to log storage

## Docker

Docker Hub
```shell script
docker pull squzy/squzy_app:v1.0.0
```

### Run locally with docker:

```shell script
docker run -p 8080:8080 squzy/squzy_app:v1.0.0
```

# Authors
- [Iurii Panarin](https://github.com/PxyUp)
- [Nikita Kharitonov](https://github.com/DreamAndDrum)

# Want to help?
Want to file a bug, contribute some code, or improve documentation? Excellent!

Add merge request with description.