# Elgo - Elo matchmaking module in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/ravsii/elgo.svg)](https://pkg.go.dev/github.com/ravsii/elgo) [![codecov](https://codecov.io/gh/ravsii/elgo/branch/main/graph/badge.svg?token=K3EM8Z6C7B)](https://codecov.io/gh/ravsii/elgo) [![Go Report Card](https://goreportcard.com/badge/github.com/ravsii/elgo)](https://goreportcard.com/report/github.com/ravsii/elgo) [![CI](https://github.com/ravsii/elgo/actions/workflows/ci.yml/badge.svg)](https://github.com/ravsii/elgo/actions/workflows/ci.yml)

Elgo is a relatively small package that provides a matchmaking pool and a simple calculator for ELO-like rating with configurable `K`-factor.

It's in the _very_ early stages of development, expect bugs and unfunushed stuff.

## Why?

The main idea is to implement some sort of a basic matchmaking tool for 3rd party apps or depelopers to use. There are plans to release it as a CLI, a Docker container and provide a server API for developers to use it as a package.

## How it works?

![How it works diagram](./docs/d2/pool.svg)
