Introduction
===
The idea behind this project is to provide a dashboard showing the status ofthe ACNodes in the system.

Deployment
===
This is really intended to be deployed as a docker container. There's an included Dockerfile for building this. It uses a multipart build to keep the final container as small as possible.

Static File Handling
===
Static files are served at /static/&lt;Version&gt;/ to ensure that a new deployment breaks caches

Configuration
===
It is possible for the dashboard to read configuration from a JSON file, however since it is intended for container deployment, it also reads from environment variables:

| Env Var | Name | Default |
| -- | -- | -- |
| MQTT_SERVER | MQTT Server to connect to | (None) |
| MQTT_CLIENTID | MQTT Client Id | ACNodeDash |
| LDAP_SERVER | LDAP Server to use for auth | (None) |
| LDAP_ENABLE | Enable LDAP authentication| false |
| LDAP_BASEDN | Base to search for users and groups in | (None) |
| LDAP_BINDDN | Bind DN to use when connecting to the LDAP server | (None) |
| LDAP_BINDPW | Password to use when connecting to LDAP | (None) |
| LDAP_USEROU | OU to search for users | ou=Users |
| LDAP_GROUPOU | OU to search for users | ou=Groups |
| LDAP_SKIPTLSVERIFY | Ignore certificate errors on LDAP connection (for testing only!) | false |

LDAP
===
This is likely to require some tweaks if it were to be used with an LDAP server
that doesn't use the same structure as LHS's:
* Users are in an OU
* Groups are in a (separate) OU
* users are linked to groups via the memberUid attribute on the group

API
===
There's an API spec in /static/api.yaml - it's worth keeping this up to date since it can be browsed and interactively played with at /swagger/

The template /templates/swagger.gohtml is basically the standard swagger standalone html,
however it has been tweaked to make it work with the cache breaking and general static file handling. 