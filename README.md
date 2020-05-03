# Squzy - opensource monitoring and alerting system

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)
[![codecov](https://codecov.io/gh/squzy/squzy/branch/develop/graph/badge.svg)](https://codecov.io/gh/squzy/squzy)
[![GolangCI](https://golangci.com/badges/github.com/squzy/golangci-lint.svg)](https://golangci.com)
[![Join the chat at https://gitter.im/squzyio/community](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/squzyio/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## About

Squzy - is a high-performance open-source monitoring and alerting system written in Golang with [Bazel](https://bazel.build/) and love.

## Contains:

### [Squzy Monitoring](https://github.com/squzy/squzy/tree/develop/apps/squzy_monitoring)

High perfomance health check system

**System Health Checks Capabilities**
1) HTTP/HTTPS
2) TCP
3) GRPC - https://github.com/grpc/grpc/blob/master/doc/health-checking.md
4) SiteMap.xml - https://www.sitemaps.org/protocol.html
5) Value from http response by selectors(https://github.com/tidwall/gjson)

### [Squzy Agents](https://github.com/squzy/squzy/tree/develop/apps/agent_client)

# Authors
- [Iurii Panarin](https://github.com/PxyUp)
- [Nikita Kharitonov](https://github.com/DreamAndDrum)

# Want to help?
Want to file a bug, contribute some code, or improve documentation? Excellent!

Add merge request with description.
