// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"storj.io/uplink"
)

// ConfigStorj depicts keys to search for within the stroj_config.json file.
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

// LoadStorjConfiguration reads and parses the JSON file that contain Storj configuration information.
func LoadStorjConfiguration(fullFileName string) ConfigStorj {

	var metric *Metric
	if useDebug {
		metric = &Metric{function: "LoadStorjConfiguration"}
		metric.start()
		defer func() {
			metric.end()
			collectedMetrics = append(collectedMetrics, metric)
		}()
	}

	var configStorj ConfigStorj
	fileHandle, err := os.Open(filepath.Clean(fullFileName))
	if err != nil {
		log.Fatal("Could not load storj config file: ", err)
	}

	jsonParser := json.NewDecoder(fileHandle)
	if err = jsonParser.Decode(&configStorj); err != nil {
		log.Fatal(err)
	}

	// Close the file handle after reading from it.
	if err = fileHandle.Close(); err != nil {
		log.Fatal(err)
	}

	// Display storj configuration read from file.
	fmt.Println("Read Storj configuration from the", fullFileName, "file.")
	fmt.Println("API Key\t\t: ", configStorj.APIKey)
	fmt.Println("Satellite	: ", configStorj.Satellite)
	fmt.Println("Bucket		: ", configStorj.Bucket)

	// Convert the upload path to standard form.
	if configStorj.UploadPath != "" {
		if configStorj.UploadPath == "/" {
			configStorj.UploadPath = ""
		} else {
			checkSlash := configStorj.UploadPath[len(configStorj.UploadPath)-1:]
			if checkSlash != "/" {
				configStorj.UploadPath = configStorj.UploadPath + "/"
			}
		}
	}

	fmt.Println("Upload Path\t: ", configStorj.UploadPath)
	fmt.Println("Serialized Access Key\t: ", configStorj.SerializedAccess)

	return configStorj
}

// ShareAccess generates and prints the shareable serialized access
// as per the restrictions provided by the user.
func ShareAccess(access *uplink.Access, configStorj ConfigStorj) {

	var metric *Metric
	if useDebug {
		metric = &Metric{function: "ShareAccess"}
		metric.start()
		defer func() {
			metric.end()
			collectedMetrics = append(collectedMetrics, metric)
		}()
	}

	allowDownload, _ := strconv.ParseBool(configStorj.AllowDownload)
	allowUpload, _ := strconv.ParseBool(configStorj.AllowUpload)
	allowList, _ := strconv.ParseBool(configStorj.AllowList)
	allowDelete, _ := strconv.ParseBool(configStorj.AllowDelete)
	notBefore, _ := time.Parse("2006-01-02_15:04:05", configStorj.NotBefore)
	notAfter, _ := time.Parse("2006-01-02_15:04:05", configStorj.NotAfter)

	permission := uplink.Permission{
		AllowDownload: allowDownload,
		AllowUpload:   allowUpload,
		AllowList:     allowList,
		AllowDelete:   allowDelete,
		NotBefore:     notBefore,
		NotAfter:      notAfter,
	}

	// Create shared access.
	sharedAccess, err := access.Share(permission)
	if err != nil {
		log.Fatal("Could not generate shared access: ", err)
	}

	// Generate restricted serialized access.
	serializedAccess, err := sharedAccess.Serialize()
	if err != nil {
		log.Fatal("Could not serialize shared access: ", err)
	}

	fmt.Println("Shareable serialized access: ", serializedAccess)
}

// ConnectToStorj reads Storj configuration from given file
// and connects to the desired Storj network.
// It then reads data property from an external file.
func ConnectToStorj(configStorj ConfigStorj, accesskey bool) (*uplink.Access, *uplink.Project) {

	var metric *Metric
	if useDebug {
		metric = &Metric{function: "ConnectToStorj"}
		metric.start()
		defer func() {
			metric.end()
			collectedMetrics = append(collectedMetrics, metric)
		}()
	}

	var access *uplink.Access
	var cfg uplink.Config

	// Configure the UserAgent
	/* For a list of valid User Agents, refer to */
	cfg.UserAgent = ""
	ctx := context.Background()
	var err error

	if accesskey {
		fmt.Println("Connecting to Storj network using Serialized access.")
		// Generate access handle using serialized access.
		access, err = uplink.ParseAccess(configStorj.SerializedAccess)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Connecting to Storj network.")
		// Generate access handle using API key, satellite url and encryption passphrase.
		access, err = cfg.RequestAccessWithPassphrase(ctx, configStorj.Satellite, configStorj.APIKey, configStorj.EncryptionPassphrase)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Open a new porject.
	project, err := cfg.OpenProject(ctx, access)
	if err != nil {
		log.Fatal(err)
	}
	defer project.Close()

	// Ensure the desired Bucket within the Project
	_, err = project.EnsureBucket(ctx, configStorj.Bucket)
	if err != nil {
		log.Fatal(err)
	}

	return access, project
}

// UploadData uploads the backup file to storj network.
func UploadData(project *uplink.Project, configStorj ConfigStorj, uploadFileName string, fileReader *os.File) {

	var metric *Metric
	if useDebug {
		metric = &Metric{function: "UploadData"}
		metric.start()
		defer func() {
			metric.end()
			collectedMetrics = append(collectedMetrics, metric)
		}()
	}

	ctx := context.Background()

	// Create an upload handle.
	upload, err := project.UploadObject(ctx, configStorj.Bucket, configStorj.UploadPath+filepath.Base(uploadFileName), nil)
	if err != nil {
		log.Fatal("Could not initiate upload : ", err)
	}
	fmt.Printf("Uploading %s to %s...\n", configStorj.UploadPath+filepath.Base(uploadFileName), configStorj.Bucket)

	// ****Add the code here to create the reader for the file to be uploaded****

	/* To directly copy the complete data to storj network, uncomment this code
	and remvove/comment the section reader code snippet.

	_, err = io.Copy(upload, fileReader)
	if err != nil {
		abortErr := upload.Abort()
		log.Fatal("Could not upload data to storj: ", err, abortErr)
	}

	*/

	// To implement uploading in parts, use the following approcach.
	// This approach creates a section reader for the file handle from the current index
	// to read the data in buffer with specified size and upload the corresponding data in sections.

	dataProcessingAndCopy(upload, fileReader)

	/*	In case you have passed a byte array(buffer) to be uploaded,
		comment the Copy function block and use the following approach.
		This approach creates a reader for 32KB section starting from the current position,
		copies the 32KB buffer data and updaes the current position.

		var lastIndex = 0
		var buf = make([]byte, 32768)

		// Loop to read the backup file in chunks and append the contents to the upload object.
		for lastIndex < int(len(dataToUpload)) {
			reader := bytes.NewBuffer(dataToUpload[lastIndex:min(lastIndex+cap(buf), len(dataToUpload))])

			_, err = io.Copy(upload, reader)

			lastIndex = lastIndex + cap(buf)
		}

	*/

	// Commit the upload after copying the complete content of the backup file to upload object.
	fmt.Println("Please wait while the upload is being committed to Storj.")
	err = upload.Commit()
	if err != nil {
		log.Fatal("Could not commit object upload : ", err)
	}

	// Close file handle after reading from it.
	if err = fileReader.Close(); err != nil {
		log.Fatal(err)
	}
}

// dataProcessingAndCopy implements the approcachof uploading data/file in parts.
// Code to modify the data to be uploaded can be added inside this function.
// By default, no modification in the uploading data has been performed.
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

/*	Uncomment this function if you are passing byte array(buffer) to the UploadData funtion.

	func min(a, b int) int {

	if a < b {
		return a
	}
	return b
}

*/
