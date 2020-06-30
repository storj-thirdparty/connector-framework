// Module to connect to a `source` instance
// and fetch data to be uploaded.

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ConfigLocalFile stores the local file path.
//****Change this strcuture to store the configurations and credentials
//     of whatever the source being used****
type ConfigLocalFile struct {
	Path string `json:"path"`
}

// LoadLocalProperty reads and parses the configuration JSON file
// that contains a local file path
// and returns it embedded in a configuration object.
//****Change the function name and add print statements as per the required configurations****
func LoadLocalProperty(fullFileName string) ConfigLocalFile {

	start := time.Now()

	var configLocalFile ConfigLocalFile

	// Open the file and generate file handle.
	fileHandle, err := os.Open(filepath.Clean(fullFileName))
	if err != nil {
		log.Fatal("Could not load influx cofig file: ", err)
	}

	// Decode and parse the JSON properties.
	//****Change the config object name here****
	jsonParser := json.NewDecoder(fileHandle)
	if err = jsonParser.Decode(&configLocalFile); err != nil {
		log.Fatal(err)
	}

	// Close the file handle after reading from it.
	if err = fileHandle.Close(); err != nil {
		log.Fatal(err)
	}

	//****Display the parsed configuration properties****
	fmt.Println("\nRead local file configuration from the ", fullFileName, " file")
	fmt.Println("File Path\t", configLocalFile.Path)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Loaded source configuration in %s with %d MiB system memory used till now.\n", time.Since(start), bToMb(m.Sys))

	return configLocalFile
}

// ConnectToLocalDisk takes the configuration object as argument
// and returns the reader of the source file to be uploaded.
//****Modify the function to connect to the required source instance
//     and return the file(s) or reader of the file****
func ConnectToLocalDisk(configLocalFile ConfigLocalFile) *os.File {

	start := time.Now()
	//****Code to connect to source and create a source instance(as per requirement)****

	//****Code fetch backup data/file(s)****

	//****Code to fetch/generate file(s)/reader****
	reader, err := os.Open(filepath.Clean(configLocalFile.Path))
	if err != nil {
		log.Fatal()
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("\nConnected to source in %s with  %d MiB system memory used till now.\n", time.Since(start), bToMb(m.Sys))

	return reader
}
