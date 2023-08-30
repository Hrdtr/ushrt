# Prerequisites

## Install golang-migrate

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Install Sqlc

```bash
# go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# Please see https://github.com/sqlc-dev/sqlc/issues/2149
# Workaround (for linux distros with snap installed):
sudo snap install sqlc
```

## Install GCC

C compiler "GCC" is required by CGO

```bash
sudo apt install build-essential
# Or:
# sudo apt install gcc
```

The build-essential package has its dependencies as gcc, g++ (GCC, but for C++), make and dpkg-dev (to build .deb packages).

# Command References

## Create New Migration:

```bash
migrate create -ext sql -dir db/migrations [name]
```

## Execute migration:

```bash
migrate -source file://db/migrations -database postgres://pguser:pgpassword@localhost:5432/starter?sslmode=disable [up/down] [n=all]
```

# Code Generation

```bash
go generate ./db
```
