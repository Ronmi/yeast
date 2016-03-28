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

## Servers

Servers is a hash table defines one or more hosts, which would contains one or more mappings, AKA `server` section in nginx.

```js
{
  "string": {"string": mapping}
}
```

# API methods

## /api/list - Lists all registered servers

This will return a `Servers`, denotes all known data.

## /api/create - create a mapping entry

By passing `name`, `path`, `upstream` and optional `custom_tags`, it will create a mapping.

The `name` can be `host` or `host:port`.

This method will return the modified `Servers` with its all paths.

## /api/modify - modify a mapping entry

By passing `name`, `path`, `new_path`, `new_upstream` and optional `new_custom_tags`, it will modify a mapping.

The `name` can be `host` or `host:port`.

This method will return the modified `Servers` with its all paths.

## /api/delete - delete a mapping entry

By passing `name` and `path`, the matching data will be deleted.

This method will return the modified `Servers` with its all paths.

**A server with zero `Mapping` means the server is deleted.**

## /api/enable - enable a mapping entry

By passing `name` and optional `path`, the matching data will be enabled.

It will enable all known settings if not passing any parameter.

This method will return the modified `Servers` with its all paths.

## /api/disable - disable a mapping entry

By passing `name` and optional `path`, the matching data will be disabled

It will disable all known settings if not passing any parameter.

This method will return the modified `Servers` with its all paths.
