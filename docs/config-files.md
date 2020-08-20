# <b>Config Files</b>

> There are two config files that contain Storj network and framework connection information. The tool is designed so you can specify a config file as part of your tooling/workflow.

## `local.json`

Inside the `./config` directory there is a `local.json` file, with following information about your framework(local file in this case) instance:

* `path`- Path to local file

## `storj_config.json`

Inside the `./config` directory a `storj_config.json` file, with Storj network configuration information in JSON format:

* `apiKey` - API Key created in Storj Satellite GUI(mandatory)
* `satelliteURL` - Storj Satellite URL(mandatory)
* `encryptionPassphrase` - Storj Encryption Passphrase(mandatory)
* `bucketName` - Name of the bucket to upload data into(mandatory)
* `uploadPath` - Path on Storj Bucket to store data (optional) or "" or "/". (mandatory)
* `serializedAccess` - Serialized access shared while uploading data used to access bucket without API Key (mandatory while using *accesskey* flag)
* `allowDownload` - Set *true* to create serialized access with restricted download (mandatory while using *share* flag)
* `allowUpload` - Set *true* to create serialized access with restricted upload (mandatory while using *share* flag)
* `allowList` - Set *true* to create serialized access with restricted list access
* `allowDelete` - Set *true* to create serialized access with restricted delete
* `notBefore` - Set time that is always before *notAfter*
* `notAfter` - Set time that is always after *notBefore*
