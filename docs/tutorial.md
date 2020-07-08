# Tutorial

> Welcome! This is the tutorial for creating your own connetor using this framework. Perform the following changes in the framework code base to create your very own connector:

## 1) Source Configuration File

* Change the name of the source config file(local.json in the sample framework) to the source name and add the required configurations and credentials to it or you can simply create your own *source* configuration file.

## 2) Source.go

Change the name of the file to *source name* and make the following amendments:

* Make required changes in the following source configuration structure to store the specified source configurations.

```
type ConfigLocalFile struct {
	Path string `json:"path"`
}
```

* Add print statements in the Load<Source>Property function to print the specified configurations.

* Add your code to connect to source and create an instance, create/fetch back-up data or files, and create a reader to the backup file/data.

## 3) Store.go

Following changes need to be made in the store.go file:

* Change the following variable and flag names in the *init()* function as per requirement and convenience.

```
var defaultLocalFile string
var defaultStorjFile string
storeCmd.Flags().BoolP("accesskey", "a", false, "Connect to storj using access key(default connection method is by using API Key).")
storeCmd.Flags().BoolP("share", "s", false, "For generating share access of the uploaded backup file.")
storeCmd.Flags().StringVarP(&defaultLocalFile, "local", "l", "././config/local.json", "full filepath contaning local file path.")
```

After making changes in the above code, you need to make the respective changes inside the *sourceStore()*(localStore in the sample code) function code as well.

* Process the upload file name to convert to a standard and less complex form as per convenience.

* Make changes in the code fragment calling the upload function as per the arguments you wish to pass. Only the *reader* and *file path/name* arguments should be changed. Single file/reader can be passed as done in the sample code. In case you have more than one file, replace the *UploadData(project, storjConfig, configLocalFile.Path, reader)* statement with the following code fragement:

```
for i := 0; i < len(FilesList); i++ {
	file := FilesList[i]
	nextcloudReader := GetReader(nextcloudClient, file)
	UploadData(project, storjConfig, file, nextcloudReader, FilesList[i])
	}
```

## 4) Storj.go

The following changes need to be made only in the upload function:

* Change the user agent to the required one. [Refer this link for valid user agents](https://github.com/storj/storj/blob/23c556ae15c5cc4735746643751d0f44a96c3e5b/satellite/rewards/partners.go).

* Change the arguments in the function definition as per the arguments passed from *store.go*.

* If reader is not passed as an argument to call the upload function, add the following code fragment to create one. Remember to close the reader after the upload is committed.

```
fileReader, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		log.Fatal(err)
	}
```
Here, filePath is the complete path of the file that needs to be uploaded.

* In case you want to implement section uploading, use the following code fragment. The corresponding code snippet has been used in the sample connector code provided.

```
func dataProcessingAndCopy(upload *uplink.Upload, fileReader *os.File) {

	var lastIndex int64
	var numOfBytesRead int
	var buf = make([]byte, 32768)
	var err1 error

	// Loop to read the backup file in chunks and append the contents to the upload object.
	for err1 != io.EOF {
		sectionReader := io.NewSectionReader(fileReader, lastIndex, int64(cap(buf)))
		numOfBytesRead, err1 = sectionReader.ReadAt(buf, 0)
		if numOfBytesRead > 0 {
			reader := bytes.NewBuffer(buf[0:numOfBytesRead])
			_, _ = io.Copy(upload, reader)
		}
		lastIndex = lastIndex + int64(numOfBytesRead)
	}
}
```

Call the above function inside the *UploadData* funciton inside *storj.go* after creating the *uplink.Upload* handle object. This approach creates a section reader for the file handle from the current index to read the data in buffer with specified size and upload the corresponding data in sections.

* For uploading a byte array(buffer), use the following code fragment. A commented block has also been provided. Uncomment the same and use it for the purpose.

```
var lastIndex = 0
var buf = make([]byte, 32768)

// Loop to read the backup file in chunks and append the contents to the upload object.
for lastIndex < int(len(dataToUpload)) {
	reader := bytes.NewBuffer(dataToUpload[lastIndex:min(lastIndex+cap(buf), len(dataToUpload))])

	_, err = io.Copy(upload, reader)

	lastIndex = lastIndex + cap(buf)
}
```

This approach creates a reader for 32KB section starting from the current position, copies the 32KB buffer data and updaes the current position. The *min* function is used to avoid referring to a null memory. Code for *min* function is as follows:

```
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

## 5) Change the connector name in *root.go* and *main.go* files.

## 6) Create a *go.mod* file for the respective connector.
