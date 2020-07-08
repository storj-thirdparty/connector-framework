# Run

> Back-up is uploaded by streaming to the Storj network.

The following flags can be used with the `store` command:

* `accesskey` - Connects to the Storj network using a serialized access key instead of an API key, satellite url and encryption passphrase.
* `shared` - Generates a restricted shareable serialized access with the restrictions specified in the Storj configuration file.
* `debug` - Prints the execution time, memory used by each function and collects the garbage memory at the end of the command execution.
Once you have built the project you can run the following:

## Get help

```
$ ./connector-framework --help
```

## Check version

```
$ ./connector-framework --version
```

## Upload back-up data to Storj

```
$ ./connector-framework store --local <path_to_local_config_file> --storj <path_to_storj_config_file>
```

## Upload back-up data to Storj bucket using Access Key

```
$ ./connector-framework store --accesskey
```

## Upload back-up data to Storj and generate a Shareable Access Key based on restrictions in `storj_config.json`

```
$ ./connector-framework store --share
```

## Upload back-up data to storj in debug mode

```
$ ./connector-framework store --debug
```
