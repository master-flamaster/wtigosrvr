# wtigosrvr

This repository provides a minimal time logger written in Go. The program
queries several public NTP servers and records their timestamps. Each log
entry is chained with an HMAC so tampering can be detected.

## Building

```
go build ./cmd/server
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
