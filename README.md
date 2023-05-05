# Elgo - Elo matchmaking module in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/ravsii/elgo.svg)](https://pkg.go.dev/github.com/ravsii/elgo) [![codecov](https://codecov.io/gh/ravsii/elgo/branch/main/graph/badge.svg?token=K3EM8Z6C7B)](https://codecov.io/gh/ravsii/elgo) [![Go Report Card](https://goreportcard.com/badge/github.com/ravsii/elgo)](https://goreportcard.com/report/github.com/ravsii/elgo)

Elgo is a relatively small package that provides a matchmaking pool and a simple calculator for ELO-like rating with configurable `K`-factor.

## How it works?

TLDR version:

![](./docs/d2/pool.svg)

## TODO list

- [ ] Add example explanation
- [ ] Add d2lang diagrams
- [ ] Add a git tag/version
- [ ] Add other pool types (non-elo based)
  - [ ] LIFO (Stack)
  - [ ] FIFO (Queue)
  - [ ] Other ...?
- [ ] Add an option to use this as a service
  - gRPC (does it even support channel-like streaming data type?)
  - as a docker container with sockets/grpc
