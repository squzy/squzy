# Squzy Agents - open source agent for grab information from host machine

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)

## About

This is small application which grab information from host machine

Which information you can get:
1. CPU load per each
2. Memory usage (used/free/total/shared)
3. Disk (used/free/total) per each disk
4. Net (bytes sent/get, package sent/get , err stat)

## Usage

For use that you should implement this [service](https://github.com/squzy/squzy_proto/blob/develop/proto/v1/squzy_agent_server.proto#L10)

Some times later we will provide our implementation

### Timeline:

1. Agent start
2. Register request (get ID)
3. SendMetrics every interval
4. ....
5. Terminate/kill signal
6. UnRegister request
7. Finish


## Environment variables

Bold is required

- **SQUZY_AGENT_SERVER_HOST** - agent server host
- SQUZY_AGENT_INTERVAL(5s) - how offen get metric
- SQUZY_AGENT_NAME - uniq name for define strict agent
