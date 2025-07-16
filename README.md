# wtigosrvr

This repository provides a minimal time logger written in Go. The program
queries several public NTP servers and records their timestamps. Each log
entry is chained with an HMAC so tampering can be detected.

## Building

```
go build ./cmd/server

# build for Linux (useful for Vultr)
GOOS=linux GOARCH=amd64 go build -o ntp-monitor ./cmd/server
```

Run with a secret key file:

```
./server -key hmac.key
```

The application periodically queries the following servers:

- time.google.com
- time.apple.com
- time.windows.com
- 0.ru.pool.ntp.org
- 1.ru.pool.ntp.org
- 2.ru.pool.ntp.org
- 3.ru.pool.ntp.org

Results are appended to `ntp.log` in JSON format.

## Deploying on Vultr

The repository includes example configuration files for provisioning a server.
`vultr_api_payload.json` contains the payload used to create an instance.
Place your API key in `vultr_config.json` and run:

```
go run vultr.go -config vultr_config.json -payload vultr_api_payload.json
```

The `readme.md` file also contains a cloud-init snippet and a macOS tray
utility example.

