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

### Execute from sources

Execute (with default configuration `dev`):
```bash
go run src/finances-service/main.go
```

Execute with a specific configuration file (example: `production`):
```bash
FINANCES_SERVICE_CONFIG=config/production.yaml go run src/finances-service/main.go
```

### Build & Execute

This outputs `./finances-service` (`.exe` on Windows) executable.

To start it with a specific configuration, prepend the `FINANCES_SERVICE_CONFIG` environment
variable assignment before calling the executable.

```bash
go build ./src/finances-service
FINANCES_SERVICE_CONFIG=config/production.yaml ./finances-service
# .\finances-service.exe on Windows
```

# Configuration

## Configuration File

By default, the service will try to load `config/dev.yaml`. By change this
behaviour, this environment variable must be set.

```bash
FINANCES_SERVICE_CONFIG=config/integration.yaml
```

## Specific Configuration Key

To set a specific configuration key, per instance `database.password`,
set an environment variable with its path in UPPERCASE, separated by underscores `_` instead
of dots `.` and prefixing it by `FINANCES_SERVICE`.

```bash
FINANCES_SERVICE_DATABASE_PASSWORD=thisISohSOsecure
```
