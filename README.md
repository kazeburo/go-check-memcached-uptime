# go-check-memcached-uptime

Mackerel check plugin for memcached uptime


# Usage

```
Usage:
  check-memcached-uptime [OPTIONS]

Application Options:
  -H, --host=     Hostname (default: localhost)
  -p, --port=     Port (default: 11211)
  -t, --timeout=  Seconds before connection times out (default: 10)
  -c, --critical= critical if uptime seconds is less than this number
  -w, --warning=  warning if uptime seconds is less than this number
  -v, --version   Show version

Help Options:
  -h, --help      Show this help message

```


Sample

```
% ./check-memcached-uptime -w 30 -c 30
memcached Uptime CRITICAL: up 0 days, 00:00:27 < 0 days, 00:00:30
% ./check-memcached-uptime -w 30 -c 30
memcached Uptime OK: up 0 days, 00:00:35
```

## Install

Please download release page or `mkr plugin install kazeburo/go-check-memcached-uptime`.

