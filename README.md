## connector-framework (uplink v1.0.5)

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


## Create your own connector

The following changes need to be made to the framework to create your own connector:

#### Source Configuration File

* Change the name of the source config file to the source name and add the required configurations and credentials to it or you can simply create your own *source* configuration file.

#### Source.go

Change the name of the file to *source name* and make the following amendments:

* Change the source configuration structure to store the specified source configurations.
* Add print statements in the Load<Source>Property function to print the specified configurations.
* Add the code to connect to source and create an instance, create/fetch back-up data or files, and create a reader to the backup file/data.

#### Store.go

Following changes need to be made in the store.go file:

* Change the corresponding variable and flag names.
* Process the upload file name to convert to a standard and less complex form, if required.
* Made changes in the code fragment calling the upload function as per the arguments you wish to pass. Only the *reader* and *file path/name* arguments should be changed.

#### Storj.go

The following changes need to be made only in the upload function:

* Change the arguments in the function definition as per the arguments passed from *store.go*.
* If reader is not passed as an argument to call the upload function, add the code fragment to create one. Remember to close the reader after the upload is committed.

#### Change the connector name in *root.go* and *main.go*.

#### Create a *go.mod* file for the respective connector.



## Documentation

For more information on runtime flags, configuration, testing, and diagrams, check out the [Detail](//github.com/storj-thirdparty/storj-framework/wiki/Home) or jump to:

* [Config Files](//github.com/storj-thirdparty/connector-framework/wiki/#config-files)
* [Run (long version)](//github.com/storj-thirdparty/connector-framework/wiki/#run)
* [Testing](//github.com/storj-thirdparty/connector-framework/wiki/#testing)
* [Flow Diagram](//github.com/storj-thirdparty/connector-framework/wiki/#flow-diagram)

