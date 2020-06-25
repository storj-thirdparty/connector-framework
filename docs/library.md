## <b>Functions</b>

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
