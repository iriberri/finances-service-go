finances-service-go
===

# Project

The back-end service behind the _Finances_ (until a better name is found) system.

# Useful commands

## Development Environment Setup

### Prerequisites

* [GO Language](https://golang.org/doc/install) development kit (>=1.10.3)
* [dep](https://golang.github.io/dep/docs/installation.html) GO dependencies tool

### Get ready to develop

```bash
git clone git@github.com:Adeynack/finances-service-go.git
cd finances-service-go
dep ensure
```

### Test

```bash
go test ./...
```

### Build & Execute

This outputs `./finances-service` (`.exe` on Windows) executable.

```bash
go build ./src/finances-service
./finances-service
# .\finances-service.exe on Windows
```
