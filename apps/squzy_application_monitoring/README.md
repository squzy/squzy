# Squzy API - open source aggregation over grpc API

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)

## About

This project provide [API](https://github.com/squzy/squzy_proto/blob/develop/proto/v1/squzy_application_monitoring.proto)

## Environment variables

Bold is required

- PORT(9095) - port for listening
- TRACING_HEADER(Squzy_transaction) - header for transfer
- **MONGO_URI** - mongo URI
- MONGO_DB(applications_monitoring) - mongo collection
- MONGO_COLLECTION(application) - mongo collection
- **SQUZY_STORAGE_HOST** - host for storage server
- SQUZY_STORAGE_TIMEOUT(5s) - timeout for storage
