# iec61850 — IEC 61850 MMS, GOOSE and SV Go binding

[![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-GPL--3.0-green.svg)](https://www.gnu.org/licenses/gpl-3.0.html)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/boeboe/iec61850)](https://pkg.go.dev/github.com/boeboe/iec61850)
[![Go Report Card](https://goreportcard.com/badge/github.com/boeboe/iec61850)](https://goreportcard.com/report/github.com/boeboe/iec61850)

Go (cgo) bindings for [libIEC61850](https://github.com/mz-automation/libiec61850) — MMS, GOOSE, and Sampled Values (SV) client/server.

## Quick start

```bash
go get -u github.com/boeboe/iec61850
```

**Requires:** libiec61850 v1.6.1 built with R-GOOSE, R-SMV, and SNTP support. See [Building libiec61850](#building-libiec61850-v161) below.

## Documentation

Detailed developer documentation:

| Document | Description |
|----------|-------------|
| [FUNCTIONS.md](FUNCTIONS.md) | API reference for all exported functions (C↔Go mapping, parameters, examples) |
| [STRUCTS.md](STRUCTS.md) | Struct and type reference (Client, Server, MMS, GOOSE, SV, configuration, timestamps) |
| [ENUMS.md](ENUMS.md) | Enums and constants (MmsType, FC, quality flags, control models, etc.) |

## Features

- **MMS** client/server, data access (read/write), discovery, data sets, reports, file services, setting groups
- **GOOSE** (IEC 61850-8-1): publish/subscribe, control blocks (GoCB), R-GOOSE
- **Sampled Values (SV)** (IEC 61850-9-2): publish/subscribe, R-SMV
- **Server:** data model, attribute updates (including timestamp + time quality), write/control handlers, TLS
- **Client:** connect (sync/async, with auth), read/write, timestamps with quality (`UtcTimeValue`), reports, TLS

## Example usage

- [Client control](test/client_control/client_control_test.go) · [Client RCB](test/client_rcb/client_rcb_test.go) · [Client read/write](test/client_rw) · [Client setting groups](test/client_sg/client_sg_test.go)
- [TLS client](test/tls_client/client_read_test.go) · [TLS server](test/tls_server/tls_server_test.go)
- [Server write access](test/server/complexModel_test.go) · [Server control](test/server/simpleIO_control_test.go) · [Server direct control + GOOSE](test/server/simpleIO_direct_control_goose_test.go)

## Building libiec61850 v1.6.1

Build the C library with R-GOOSE, R-SMV, and SNTP before using this package.

**Prerequisites:** mbedtls (e.g. `brew install mbedtls` on macOS, `apt-get install libmbedtls-dev` on Debian/Ubuntu).

```bash
./scripts/rebuild_libraries.sh
go test -v -run TestLibraryVersion   # verify build
```


## License

**This project is licensed under the GNU General Public License v3.0 (GPL-3.0).**

This is required because the Go binding links against **[libIEC61850](https://github.com/mz-automation/libiec61850)** (the official open-source C library for IEC 61850 protocols by MZ Automation), which is distributed under [GPL-3.0](https://www.gnu.org/licenses/gpl-3.0.html). Any use or distribution of this binding is therefore subject to the same GPL-3.0 terms. See [LICENSE](./LICENSE) in this repository for the full license text.

