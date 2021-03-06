# libgrin

![CI](https://github.com/blockcypher/libgrin/workflows/Continuous%20Integration/badge.svg)
[![ISC License](http://img.shields.io/badge/license-apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/blockcypher/libgrin/v5.svg)](https://pkg.go.dev/github.com/blockcypher/libgrin/v5)

libgrin is a golang Grin library for building [Grin](https://github.com/mimblewimble/grin) applications.
Currently, it contains the basic consensus parameters, chain, slate structures and proof of work verification code.

The `client` package contains wrappers around the Grin node foreign/owner API and the wallet foreign/owner API using libgrin.

## Requirements

[Go](http://golang.org) 1.14 or newer.

## Installation

### Linux/macOS

- Install Go according to the installation instructions [here](http://golang.org/doc/install).

- Run the following commands to obtain libgrin, all dependencies, and build it:

```bash
cd $GOPATH/src/github.com/blockcypher/libgrin
go get ./...
go build ./...
```

- libgrin (and utilities) is now built and can be used in your program.

## Issue Tracker

The [integrated github issue tracker](https://github.com/blockcypher/libgrin/issues)
is used for this project.

## Contributing

Any contribution is more than welcome. Currently libgrin has:

- Minimal API types.
- Minimal Core package which includes PoW verification.
- Client package containing wrappers around Node API (Foreign and Owner) and wallet API (Foreign and Owner).
- Minimal keychain types
- Minimal libwallet
- Minimal p2p types
- Minimal pool types

Missing parts:

- Complete wallet
- Complete keychain
- All the rest

This project tries to follow as much as possible the official [Grin node code structure](https://github.com/mimblewimble/grin) and [Grin-Wallet code structure](https://github.com/mimblewimble/grin-wallet).

## License

libgrin is licensed under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) license.
