## Flow Diagram

![](https://github.com/storj-thirdparty/storj-framework/blob/master/README.assets/arch.drawio.png)



## Config Files

There are two config files that contain Storj network and framework connection information. The tool is designed so you can specify a config file as part of your tooling/workflow.

##### `local.json`

Inside the `./config` directory there is a `local.json` file, with following information about your framework(local file in this case) instance:

* `path`- Path to local file

##### `storj_config.json`

Inside the `./config` directory a `storj_config.json` file, with Storj network configuration information in JSON format:

* `apiKey` - API Key created in Storj Satellite GUI(mandatory)
* `satelliteURL` - Storj Satellite URL(mandatory)
* `encryptionPassphrase` - Storj Encryption Passphrase(mandatory)
* `bucketName` - Name of the bucket to upload data into(mandatory)
* `uploadPath` - Path on Storj Bucket to store data (optional) or "/" (mandatory)
* `serializedAccess` - Serialized access shared while uploading data used to access bucket without API Key (mandatory while using *accesskey* flag)
* `allowDownload` - Set *true* to create serialized access with restricted download (mandatory while using *share* flag)
* `allowUpload` - Set *true* to create serialized access with restricted upload (mandatory while using *share* flag)
* `allowList` - Set *true* to create serialized access with restricted list access
* `allowDelete` - Set *true* to create serialized access with restricted delete
* `notBefore` - Set time that is always before *notAfter*
* `notAfter` - Set time that is always after *notBefore*



## Run

Back-up is uploaded by streaming to the Storj network.

The following flags can be used with the `store` command:

* `accesskey` - Connects to the Storj network using a serialized access key instead of an API key, satellite url and encryption passphrase.
* `shared` - Generates a restricted shareable serialized access with the restrictions specified in the Storj configuration file.
* `debug` - Prints the execution time, memory used by each function and collects the garbage memory at the end of the command execution.

Once you have built the project you can run the following:

##### Get help

```
$ ./connector-framework --help
```

##### Check version

```
$ ./connector-framework --version
```

##### Upload back-up data to Storj

```
$ ./connector-framework store --local <path_to_local_config_file> --storj <path_to_storj_config_file>
```

##### Upload back-up data to Storj bucket using Access Key

```
$ ./connector-framework store --accesskey
```

##### Upload back-up data to Storj and generate a Shareable Access Key based on restrictions in `storj_config.json`

```
$ ./connector-framework store --share
```

##### Upload back-up data to storj in debug mode

```
$ ./connector-framework store --debug
```


## Testing

The project has been tested on the following operating systems:

```
* Windows
	* Version: 10 Pro
	* Processor: Intel(R) Core(TM) i3-5005U CPU @ 2.00GHz 2.00GHz

* macOS Catalina
	* Version: 10.15.4
	* Processor: 2.5 GHz Dual-Core Intel Core i5

* ubuntu
	* Version: 16.04 LTS
	* Processor: AMD A6-7310 APU with AMD Radeon R4 Graphics × 4
```



## Functions

#### ConnectToLocalDisk

```
func ConnectToLocalDisk(configLocalFile ConfigLocalFile) *os.File
```

ConnectToLocalDisk takes the configuration object as argument and returns the reader of the source file to be uploaded. Function name and implementation can be changed to connect to the required source instance and return the required handle/reader.

#### ConnectToStorj

```
func ConnectToStorj(fullFileName string, configStorj ConfigStorj, accesskey bool) (*uplink.Access, *uplink.Project)
```

ConnectToStorj reads Storj configuration from given file and connects to the desired Storj network. It then reads data property from an external file.

#### ShareAccess

```
func ShareAccess(access *uplink.Access, configStorj ConfigStorj)
```

ShareAccess generates and prints the shareable serialized access as per the restrictions provided by the user.
 
#### UploadData

```
func UploadData(project *uplink.Project, configStorj ConfigStorj, uploadFileName string, fileReader *os.File)
```

UploadData uploads the backup file to storj network. Parameters can be changed as per the requirement. If reader/handle is not passed as an argument to call the function, add the corresponding code snippet to create one. Remember to close the reader after the upload is committed.




## Types

#### ConfigLocalFile

```
type ConfigLocalFile struct {
	Path string `json:"path"`
}
```

ConfigLocalFile stores the local file path. Change the strcuture name and definition to store the configurations and credentials of whatever the source being used.

#### ConfigStorj

```
type ConfigStorj struct {
	APIKey               string `json:"apikey"`
	Satellite            string `json:"satellite"`
	Bucket               string `json:"bucket"`
	UploadPath           string `json:"uploadPath"`
	EncryptionPassphrase string `json:"encryptionpassphrase"`
	SerializedAccess     string `json:"serializedAccess"`
	AllowDownload        string `json:"allowDownload"`
	AllowUpload          string `json:"allowUpload"`
	AllowList            string `json:"allowList"`
	AllowDelete          string `json:"allowDelete"`
	NotBefore            string `json:"notBefore"`
	NotAfter             string `json:"notAfter"`
}
```

ConfigStorj depicts keys to search for within the stroj_config.json file.


## Tutorial

The following cis the tutorial to create your own connector:

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
* In case you want to implement section uploading, the corresponding code snippet has been provided in a commented block for the same. Uncomment the same and use it for the purpose.
* For uplaoding a byte array(buffer), a commented block has been provided for that also. Uncomment the same and use it for the purpose.

#### Change the connector name in *root.go* and *main.go*.

#### Create a *go.mod* file for the respective connector.
