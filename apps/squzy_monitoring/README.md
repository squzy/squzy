# Squzy Monitoring - open source health check service

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)

## About

System allows monitoring resources with different type of check. It uses MONGO for save state of the application and restore

### System Health Checks Capabilities
1) HTTP/HTTPS
2) TCP
3) GRPC - https://github.com/grpc/grpc/blob/master/doc/health-checking.md
4) SiteMap.xml - https://www.sitemaps.org/protocol.html
5) Value from http response by selectors(https://github.com/tidwall/gjson)
6) SSL Expiration - monitoring when SSL cert is over

# Usage

## API

[**GRPC API**](https://github.com/squzy/squzy_proto/blob/master/proto/v1/squzy_monitoring.proto) 

## Storage

This is entity for save results from squzy monitoring

You can implement storage by your self via this grpc method [API](https://github.com/squzy/squzy_proto/blob/feat/api-1/proto/v1/squzy_storage.proto#L19)

By default squzy monitoring will send **success checks in stdout**, **errors in stderr**


# Examples of call from [BloomRPC](https://github.com/uw-labs/bloomrpc)

### Http/Https check:

Usually that check used for monitoring web sites

```shell script
{
  "interval": 10, - 10 second interval
  "timeout": 5, - // default timeout is 10 sec
  "http": {
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
  "timeout": 5, - // default timeout is 10 sec
  "tcp": {
    "host": "localhost", - host
    "port": 6345 - port
  },
}
```

### SSL Expiration check:

Check can be used for validate SSL cert

```shell script
{
  "interval": 10, - 10 second interval
  "timeout": 5, - // default timeout is 10 sec
  "ssl_expiration": {
    "host": "localhost", - host
    "port": 6345 - port
  },
}
```

### SiteMap check:

**Supports redirects!**

**Every route should return 200**

That check good usage when you have critical URL in sitemap, if any of URL throw error check will be failed

```shell script
{
  "interval": 10,
  "timeout": 5, - // default timeout is 10 sec
  "sitemap": {
    "url": "https://www.sitemaps.org/sitemap.xml", - url of sitemap (https://www.sitemaps.org/sitemap.xml)
    "concurrency": 5 - parallel 5 request  
  },
}
```

### GRPC check:

Check better to use for internal testing of API services

```shell script
{
  "interval": 10,
  "timeout": 5, - // default timeout is 10 sec
  "grpc": {
    "service": "Check", - service name
    "host": "localhost", - host
    "port": 9090 - port
  },
}
```

### Value monitoring from Http json response (v1.3.0+)

Monitoring specific value from http request by json selector

Valid selectors you can find here: https://github.com/tidwall/gjson

Support type: https://github.com/squzy/squzy_proto/blob/master/proto/v1/server.proto#L84
    

```shell script
{
  "interval": 10,
  "timeout": 5, - // default timeout is 10 sec
  "httpValue": {
      "method": "GET",
      "url": "https://api.exchangeratesapi.io/latest?base=USD",
      "headers": {
        "custom": "yes",
      },
      "selectors": [
        {
          "type": 4,
          "path": "rates.RUB"
        }
      ]
    }
}
```

## Environment variables

Bold is required

- PORT(9090) - on with port run squzy
- SQUZY_STORAGE_HOST - log storage host(example *localhost:9090*)
- SQUZY_STORAGE_TIMEOUT - timeout for connect to log storage
- **MONGO_URI** - mongo url for save data
- MONGO_DB(squzy_monitoring) - mongo db name
- MONGO_COLLECTION(schedulers) - in which collection we should save data

## Docker

[HUB](https://hub.docker.com/repository/docker/squzy/squzy_monitoring)

For current develop branch use tag: **latest**

Docker Hub

```shell script
docker pull squzy/squzy_monitoring:v1.6.0
```

### Run locally with docker:

```shell script
docker run -p 9090:9090 squzy/squzy_monitoring:v1.6.0
```

# Want to help?
Want to file a bug, contribute some code, or improve documentation? Excellent!

Add merge request with description.
