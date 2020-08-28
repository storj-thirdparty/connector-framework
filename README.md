## connector-framework (uplink v1.0.5)

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/253d84109b174697b8453e81d8998073)](https://app.codacy.com/gh/storj-thirdparty/connector-framework?utm_source=github.com&utm_medium=referral&utm_content=storj-thirdparty/connector-framework&utm_campaign=Badge_Grade_Dashboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/storj-thirdparty/connector-framework)](https://goreportcard.com/report/github.com/storj-thirdparty/connector-framework)
![Cloud Build](https://storage.googleapis.com/storj-utropic-services-badges/builds/connector-framework/branches/master.svg)

## Overview

The framework Connector is a generic connector that can be used to take backup from the specified source and upload the backup files on Storj network. Sample connector to local disk is provided.

```bash
Usage:
  connector-framework [command] <flags>

Available Commands:
  help        Help about any command
  store       Command to upload data to a Storj V3 network
  version     Prints the version of the tool
```

`store` - Connect to the specified(default: `local.json`). Back-up data are generated using tooling provided by framework then uploaded to the Storj network. Connect to a Storj v3 network using the access specified in the Storj configuration file (default: `storj_config.json`).


Sample configuration files are provided in the `./config` folder.

## Requirements and Install

To build from scratch, [install the latest Go](https://golang.org/doc/install#install).

> Note: Ensure go modules are enabled (GO111MODULE=on)

#### Option #1: clone this repo (most common)

To clone the repo

```
git clone https://github.com/storj-thirdparty/connector-framework.git
```

Then, build the project using the following:

```
cd connector-framework
go build
```

#### Option #2:  ``go get`` into your gopath

To download the project inside your GOPATH use the following command:

```
go get github.com/storj-thirdparty/connector-framework
```

> Note: For reference, connector-local to backup a local file is made and following commands can be used to test the same.

## Run (short version)

Once you have built the project run the following commands as per your requirement:

##### Get help

```
$ ./connector-framework --help
```

##### Check version

```
$ ./connector-framework --version
```

##### Create backup from framework and upload to Storj

```
$ ./connector-framework store
```

## Documentation

For more information on runtime flags, configuration, testing, and diagrams, check out the [Detail](//github.com/storj-thirdparty/storj-framework/wiki/Home) or jump to:

* [Config Files](//github.com/storj-thirdparty/connector-framework/wiki/#config-files)
* [Run (long version)](//github.com/storj-thirdparty/connector-framework/wiki/#run)
* [Testing](//github.com/storj-thirdparty/connector-framework/wiki/#testing)
* [Flow Diagram](//github.com/storj-thirdparty/connector-framework/wiki/#flow-diagram)
* [Functions](//github.com/storj-thirdparty/connector-framework/wiki/#funcitons)
* [Types](//github.com/storj-thirdparty/connector-framework/wiki/#types)
* [Tutorial](//github.com/storj-thirdparty/connector-framework/wiki/#tutorial)
