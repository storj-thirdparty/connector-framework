// Module to connect to a `source` instance
// and fetch data to be uploaded.

package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	var metric *Metric
	if useDebug {
		metric = &Metric{function: "LoadLocalProperty"}
		metric.start()
		defer func() {
			metric.end()
			collectedMetrics = append(collectedMetrics, metric)
		}()
		/*var m runtime.MemStats
		runtime.ReadMemStats(&m)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Printf("LoadLacalProperty\tStart\tCurrent RAM usage: %d MiB\n", bToMb(m.HeapInuse)+bToMb(m.StackInuse))*/
	}

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
	fmt.Println("Read local file configuration from the", fullFileName, "file.")
	fmt.Println("File Path\t: ", configLocalFile.Path)

	/*if useDebug {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Printf("LoadLacalProperty\tEnd\tCurrent RAM usage: %d MiB\n\n", bToMb(m.HeapInuse)+bToMb(m.StackInuse))
	}*/

	return configLocalFile
}

// ConnectToLocalDisk takes the configuration object as argument
// and returns the reader of the source file to be uploaded.
//****Modify the function to connect to the required source instance
//     and return the file(s) or reader of the file****
func ConnectToLocalDisk(configLocalFile ConfigLocalFile) *os.File {

	var metric *Metric
	if useDebug {
		metric = &Metric{function: "ConnectToLocalDisk"}
		metric.start()
		defer func() {
			metric.end()
			collectedMetrics = append(collectedMetrics, metric)
		}()
	}
	//****Code to connect to source and create a source instance(as per requirement)****

	//****Code fetch backup data/file(s)****

	//****Code to fetch/generate file(s)/reader****
	reader, err := os.Open(filepath.Clean(configLocalFile.Path))
	if err != nil {
		log.Fatal()
	}

	return reader
}
