# GoldSrc HLDS Input Plugin

The `hlds` plugin gather metrics from GoldSrc HLDS Server.

## Configuration

```toml @sample.conf
# Fetch metrics from a HLDS
[[inputs.hlds]]
  ## Specify servers using the following format:
  ##    servers = [
  ##      ["ip1:port1", "rcon_password1", "sv1"],
  ##      ["ip2:port2", "rcon_password2", "sv2"],
  ##    ]
  #
  ## If no servers are specified, no data will be collected
  servers = []
```

## Metrics

The plugin retrieves the output of the `stats` command that is executed via
rcon.

If no servers are specified, no data will be collected

- hlds
  - tags:
    - host
    - svname
  - fields:
    - cpu (float)
    - net_in (float)
    - net_out (float)
    - uptime_minutes (float)
    - users (float)
    - fps (float)
    - players (float)
