# Squzy Incident server

[![version](https://img.shields.io/github/v/release/squzy/squzy.svg)](https://github.com/squzy/squzy)

## About

Provide possability for handaling users incident from storage

## API

[**GRPC API**](https://github.com/squzy/squzy_proto/blob/master/proto/v1/squzy_incident_server.proto#L11) 

## Environment variables

Bold is required

- PORT(9097) - on with port run squzy_incident
- *MONGO_URI* - mongo URI for connect
- MONGO_DB(incident_manager) - mongo DB for connect
- MONGO_COLLECTION(rules) - collection name
- *STORAGE_HOST* - squzy storage host
