# GoldSrc HLDS Input Plugin

The `hlds` plugin gather metrics from GoldSrc HLDS (Half-life dedicated server).

## Global configuration options <!-- @/docs/includes/plugin_config.md -->

In addition to the plugin-specific configuration settings, plugins support
additional global and plugin configuration settings. These settings are used to
modify metrics, tags, and field or create aliases and configure ordering, etc.
See the [CONFIGURATION.md][CONFIGURATION.md] for more details.

## Configuration

```toml @sample.conf
# Fetch metrics from a GoldSrc HLDS
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
    - svid
  - fields:
    - cpu (float)
    - net_in (float)
    - net_out (float)
    - uptime_minutes (float)
    - users (float)
    - fps (float)
    - players (float)

## Example Output