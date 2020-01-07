# libgrin

![](https://github.com/blockcypher/libgrin/workflows/Continuous%20Integration/badge.svg)
[![ISC License](http://img.shields.io/badge/license-apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/blockcypher/libgrin)

libgrin is a golang Grin library for building [Grin](https://github.com/mimblewimble/grin) applications.
Currently, it contains the basic consensus parameters, chain, slate structures and proof of work verification code. The `example` package contains examples of wrappers around the Grin node API and the wallet owner API using libgrin.

## Requirements

[Go](http://golang.org) 1.11 or newer.

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
- Example package containing wrappers around Node API and wallet owner API.
- Minimal keychain types
- Minimal libwallet

Missing parts:

- Complete wallet
- Complete keychain
- All the rest

This project tries to follow as much as possible the official [Grin node code structure](https://github.com/mimblewimble/grin) and [Grin-Wallet code structure](https://github.com/mimblewimble/grin-wallet).

## License

libgrin is licensed under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) license.
