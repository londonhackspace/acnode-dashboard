Introduction
===
The idea behind this project is to provide a dashboard showing the status ofthe ACNodes in the system.

Deployment
===
This is really intended to be deployed as a docker container. There's an included Dockerfile for building this. It uses a multipart build to keep the final container as small as possible.

Configuration
===
It is possible for the dashboard to read configuration from a JSON file, however since it is intended for container deployment, it also reads from environment variables:

| Env Var | Name | Default |
| -- | -- | -- |
| MQTT_SERVER | MQTT Server to connect to | (None) |
| MQTT_CLIENTID | MQTT Client Id | ACNodeDash |
| ACSERVER_URL | ACServer URL | https://acserver.london.hackspace.org.uk |
| ACSERVER_APIKEY | ACServer API Key | (None) |
| LDAP_SERVER | LDAP Server to use for auth | (None) |
| LDAP_ENABLE | Enable LDAP authentication| false |
| LDAP_BASEDN | Base to search for users and groups in | (None) |
| LDAP_BINDDN | Bind DN to use when connecting to the LDAP server | (None) |
| LDAP_BINDPW | Password to use when connecting to LDAP | (None) |
| LDAP_USEROU | OU to search for users | ou=Users |
| LDAP_GROUPOU | OU to search for users | ou=Groups |
| LDAP_SKIPTLSVERIFY | Ignore certificate errors on LDAP connection (for testing only!) | false |
| REDIS_ENABLE | Enable Redis persistence | false |
| REDIS_SERVER | Redis server to connect to | (None) |
| ADMIN_GROUPS | Comma separated list of admin user groups | (None) |
| LOG_FMT_JSON | Boolean, should logs be emitted as JSON? | false |

LDAP
===
This is likely to require some tweaks if it were to be used with an LDAP server
that doesn't use the same structure as LHS's:
* Users are in an OU
* Groups are in a (separate) OU
* users are linked to groups via the memberUid attribute on the group

Redis
====
Redis serves as a caching layer and persistence store. It is an in-memory cache, but periodically saves
state to disk, so can also be used for longer term persistence. 

If enabled, it is used as a session store, user store, and general
data persistence store for node data.

Users
====
The code supports multiple authentication providers. It is envisaged that LDAP will be the
main one in use, however there are a few cases when another source of users may be preferred:
* development
* Use somewhere other than LHS
* Storing machine accounts for API use

There is a utility in [bootstrapper](bootstrapper) which creates a
default "admin" user, for bootstrapping new systems. It uses the same configuration environment variables and/or file as 
the main dashboard process. It will create a user in the Redis provider, though
it could easily be extended to support other writable providers.

API
===
There's an API spec in /static/api.yaml - it's worth keeping this up to date since it can be browsed and interactively played with at /swagger/

The template /templates/swagger.gohtml is basically the standard swagger standalone html,
however it has been tweaked to make it work with the cache breaking and general static file handling. 

Frontend
===
The server will serve out of frontend/dist if it exists, otherwise it 
serves out of static/ - this is so when you are developing, it will use
a UI build if you have one. 

The frontend is a React app in frontend/ - to work with it:
* `npm run serve` - serves a development server, with live updates
* `npm run build` - builds a production build to dist/
* `npm run builddev` - builds a development build to dist/

The development server runs on port 3000 and forwards /api requests
to http://localhost:8080/api - this is assumed to be your local API server.