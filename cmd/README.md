# Squzy CLI

Squzy CLI allows user to interact with application by command line interface.

## storage

storage initialize squzy_storage application with provided variables.

With CLI variables can be setup by configfile, env variables or combination of these approaches.
Application will try to read env from configfile at first, then from os env.

Variables: 
- PORT(9090) - on with port run squzy_storage
- **DB_HOST** - postgresSQL host
- **DB_PORT** - postgresSQL port
- **DB_NAME** - postgresSQL name
- **DB_USER** - postgresSQL user
- **DB_PASSWORD** - postgresSQL password
- ENV_ENABLE_INCIDENT(false) - true if we want to have incidents
- ENV_INCIDENT_SERVER_HOST - required, when ENV_ENABLE_INCIDENT=true
- DB_LOGS(false) - provide logs for DB

Usage:
```
squzy storage --config path_to_config_file
```