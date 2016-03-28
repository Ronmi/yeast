# Data structure

## Mapping

A mapping defines a path under a host, AKA `location` section in nginx.

```js
{
  "custom_tags": "string", // custom nginx settings
  "enabled": bool,         // is this enabled
  "upstream": "string"     // where to proxy the traffic, the "proxy_pass" in nginx
}
```

## Server

A server defines a host, which would contains one or more mappings, AKA `server` section i nginx.

```js
{
  "name": "string",             // host name, AKA "server_name" in nginx
  "paths": {"string": mapping}, // path mappings
}
```

# API methods

## /api/list - Lists all registered servers

This will return an array of `Server`s, denotes all known data.

## /api/set - set a mapping entry

By passing `name`, `path`, `upstream` and optional `custom_tags`, it will create/overwrite a mapping.

The `name` can be `host` or `host:port`.

This method will return the added/modified `Server` data.

## /api/unset - delete a mapping entry

By passing `name` and `path`, the matching data will be deleted.

## /api/enable - enable a mapping entry

By passing `name` and optional `path`, the matching data will be enabled.

It will enable all known settings if not passing any parameter.

This method will return the modified `Server` data.

## /api/disable - disable a mapping entry

By passing `name` and optional `path`, the matching data will be disabled

It will disable all known settings if not passing any parameter.

This method will return the modified `Server` data.
