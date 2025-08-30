# Transmission Lock (TrLock)

Auto-blocker on invalid Transmission peer clients focusing on
1. Performance
2. Lightweight
3. Security

For server-side high-volume torrenting

Planned features:
1. Auto-identifing invalid peers based on config. can support client-name, ip-range, and geolocation
2. Fast and reliable peer clients lookup and verification
3. Block by IPTable (Linux, WIP) or PF (BSD)
4. Block by Transmission blocklist (only supported if using Transmission 4.1.0-2 beta version)

## Build

Build
```
make tools
make all
```

## Config

Config file path is defined by environment variable `TRLOCK_CONFIGFILE`. Config file has JSON format

An example config file:
```
{
    "host": "127.0.0.1", // Transmission daemon rpc address
    "port": "9091", // Transmission daemon rpc port
    "allowlist": {}, // Allowlist for strict allow mode
    "blocklist": {
        "client": [ // Block by client (lowercase string)
            "bad_client_name1",
            "bad_client_name2"
        ]
    },
    "strict_allow_enabled": false, // Turn on/off strict allow mode
    "pf_enabled": true, // Turn on/off block by PF (BSD only)
    "blocklist_enabled": true, // Turn on/off block by Transmission blocklist
    "blocklist_path": "./tr_blocklist", // Transmission blocklist location
    "interval": "10m", // Auto-update interval
    "reset_interval": "24h" // Auto-reset interval (clean all blocked clients)
}
```

Environment variables:
1. `TRLOCK_CONFIGFILE`: config file path
2. `TRLOCK_LOGFILE`: log file path
3. `TRLOCK_LOGDEBUG`: turn on/off debug-level logs
