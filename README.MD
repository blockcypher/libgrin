# libgrin

[![Build Status](https://dev.azure.com/blockcypher/libgrin/_apis/build/status/blockcypher.libgrin?branchName=master)](https://dev.azure.com/blockcypher/libgrin/_build/latest?definitionId=1&branchName=master)
[![ISC License](http://img.shields.io/badge/license-apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/blockcypher/libgrin)

libgrin is a golang Grin library for building Grin applications.
Currently, it contains the basic consensus parameters and proof of work verification code. The wrapper package is an example of wrappers around the Grin node, wallet and owner APIs using libgrin.

## Requirements

[Go](http://golang.org) 1.11 or newer.

## Installation

### Linux/macOS

- Install Go according to the installation instructions [here](http://golang.org/doc/install).

- Run the following commands to obtain libgrin, all dependencies, and install it:

```bash
cd $GOPATH/src/github.com/blockcypher/libgrin
go get .
```

- libgrin (and utilities) will now be installed in ```$GOPATH/bin```.  If you did
  not already add the bin directory to your system path during Go installation,
  we recommend you do so now.

## Issue Tracker

The [integrated github issue tracker](https://github.com/blockcypher/libgrin/issues)
is used for this project.

## Contributing

Any contribution is more than welcome. Currently libgrin has:

- Minimal core package
- Minimal api types
- Minimal wallet types

Missing parts:

- Wallet
- Keychain
- All the rest

This project tries to follow as much as possible the [official Grin project code structure](https://github.com/mimblewimble/grin).

## License

libgrin is licensed under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) license.