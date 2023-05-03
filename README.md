# Elgo - Elo matchmaking module in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/ravsii/elgo.svg)](https://pkg.go.dev/github.com/ravsii/elgo) [![codecov](https://codecov.io/gh/ravsii/elgo/branch/main/graph/badge.svg?token=K3EM8Z6C7B)](https://codecov.io/gh/ravsii/elgo) [![Go Report Card](https://goreportcard.com/badge/github.com/ravsii/elgo)](https://goreportcard.com/report/github.com/ravsii/elgo)

Elgo is a relatively small package that provides a matchmaking pool and a simple calculator for ELO-like rating with configurable `K`-factor.

## Why?

The main idea is to implement some sort of a basic matchmaking tool for 3rd party apps to use. There are plans to release it as a CLI, a Docker container and provide a server API for developers to use it as a package.

## How it works?

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
