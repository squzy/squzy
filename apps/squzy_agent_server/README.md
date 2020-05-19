# Squzy Agents Server - open source agent for manage squzy agent

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)

## About

This storage implement this [API](https://github.com/squzy/squzy_proto/blob/master/proto/v1/squzy_agent_server.proto#L10)

## Environment variables

Bold is required

- **MONGO_URI** - mongo URI
- PORT(9091) - port for listening
- SQUZY_STORAGE_HOST - host for storage
- SQUZY_STORAGE_TIMEOUT(5s) - timeout for storage
- MONGO_DB(squzy_agent) - db for mongo
- MONGO_COLLECTION(agents) - collection mongo
