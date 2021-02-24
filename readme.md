# traefik-jwt-headers-plugin ![Build](https://github.com/lion7/traefik-jwt-headers-plugin/actions/workflows/main.yml/badge.svg)
Traefik middleware plugin which forwards JWT claims as request headers

## Installation
The plugin needs to be configured in the Traefik static configuration before it can be used.

### Installation with Helm
The following snippet can be used as an example for the `values.yaml` file:
```values.yaml
pilot:
  enabled: true
  token: xxxxx-xxxx-xxxx

experimental:
  plugins:
    enabled: true

additionalArguments:
- --experimental.plugins.traefik-jwt-headers-plugin.modulename=github.com/lion7/traefik-jwt-headers-plugin
- --experimental.plugins.traefik-jwt-headers-plugin.version=v0.0.3
```

### Installation via command line
```
traefik \
  --experimental.pilot.token=xxxx-xxxx-xxx \
  --experimental.plugins.traefik-jwt-headers-plugin.moduleName=github.com/lion7/traefik-jwt-headers-plugin \
  --experimental.plugins.traefik-jwt-headers-plugin.version=v0.0.3
```

## Configuration
You can decide to limit the forwarded claims/headers to a given list with the `claims` option.

Each claim can be set to:

- `keep` to keep the value
- `drop` to drop the value

The `defaultMode` for `claims` is `drop`.

## Example configuration

### Kubernetes
``` tab="File (Kubernetes)"
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: my-traefik-jwt-headers-plugin
spec:
  plugin:
    traefik-jwt-headers-plugin:
      defaultMode: drop
      claims:
        sub: keep
        mysecret: drop
        user.name: keep
        organization.name: keep
```

### File (TOML)
```toml tab="File (TOML)"
[http]
  [http.middlewares]
    [http.middlewares.my-traefik-jwt-headers-plugin]
      [http.middlewares.my-traefik-jwt-headers-plugin.plugin]
        [http.middlewares.my-traefik-jwt-headers-plugin.plugin.traefik-jwt-headers-plugin]
          defaultMode = "drop"
          [http.middlewares.my-traefik-jwt-headers-plugin.plugin.traefik-jwt-headers-plugin.claims]
            sub = "keep"
            mysecret = "drop"
            user.name = "keep"
            organization.name = "keep"
```

### File (YAML)
```yaml tab="File (YAML)"
http:
    middlewares:
        my-traefik-jwt-headers-plugin:
            plugin:
                traefik-jwt-headers-plugin:
                    defaultMode: drop
                    claims:
                        sub: keep
                        mysecret: drop
                        user.name: keep
                        organization.name: keep
```

## License
This software is released under the Apache 2.0 License
