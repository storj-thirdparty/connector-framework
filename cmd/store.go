package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// storeCmd represents the store command.
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Command to upload data to storjV3 network.",
	Long:  `Command to connect and uplaod local file to given Storj Bucket.`, // ****Change the description here****
	Run:   localStore,                                                        //****Change the name by replacing `local` with the name of the source being used****
}

func init() {

	// Setup the store command with its flags.
	rootCmd.AddCommand(storeCmd)
	var defaultLocalFile string //****Change the variable name to store the default configuration file location****
	var defaultStorjFile string
	storeCmd.Flags().BoolP("accesskey", "a", false, "Connect to storj using access key(default connection method is by using API Key).")
	storeCmd.Flags().BoolP("share", "s", false, "For generating share access of the uploaded backup file.")
	storeCmd.Flags().BoolP("debug", "d", false, "For code instrumentation and profiling.")
	storeCmd.Flags().StringVarP(&defaultLocalFile, "local", "l", "././config/local.json", "full filepath contaning local file path.") //****Change the flag name and its description****
	storeCmd.Flags().StringVarP(&defaultStorjFile, "storj", "u", "././config/storj_config.json", "full filepath contaning Storj V3 configuration.")
}

var useDebug bool
var collectedMetrics []*Metric
var start time.Time

func localStore(cmd *cobra.Command, args []string) {

	// Process arguments from the CLI.
	localConfigFilePath, _ := cmd.Flags().GetString("local") //****Change the command argument here****
	fullFileNameStorj, _ := cmd.Flags().GetString("storj")
	useAccessKey, _ := cmd.Flags().GetBool("accesskey")
	useAccessShare, _ := cmd.Flags().GetBool("share")
	useDebug, _ = cmd.Flags().GetBool("debug")
	defer func() {
		if useDebug {
			err := saveCollectedMetrics(collectedMetrics)
			if err != nil {
				fmt.Printf("failed to save metrcis %s", err)
			}
		}
	}()

	// Read local file configuration from an external file and create a configuration object.
	//****Change the statement as per the `source` code Function****
	configLocalFile := LoadLocalProperty(localConfigFilePath)

	// Read storj network configurations from and external file and create a storj configuration object.
	storjConfig := LoadStorjConfiguration(fullFileNameStorj)

	// Connect to storj network using the specified credentials.
	access, project := ConnectToStorj(storjConfig, useAccessKey)

	// Retrieve the reader of the specified file.
	//****This will store the file(s)/reader to be uplaoded****
	reader := ConnectToLocalDisk(configLocalFile)

	fmt.Printf("Initiating back-up.\n")
	// Upload the desired file to desired Storj bucket.
	//****Change this code fragment by adding a loop if more than one file are to be uploaded
	//    and also process the file name to be uplaoded to a standard form(if required)****
	UploadData(project, storjConfig, configLocalFile.Path, reader)
	fmt.Printf("Back-up complete.\n\n")

	// Create restricted shareable serialized access if share is provided as argument.
	if useAccessShare {
		ShareAccess(access, storjConfig)
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
