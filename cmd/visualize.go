package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type visMetric struct {
	Metric
	RunUUID string
}

const metricsFlag = "metrics"

var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Command to visualize performance metrics.",
	Long:  `Command to visualize locally stored performance metrics.`,
	Run:   visualizeMetrics,
}

func init() {
	// Setup the store command with its flags.
	rootCmd.AddCommand(visualizeCmd)
	var metricsFolder string
	visualizeCmd.Flags().StringVarP(&metricsFolder, metricsFlag, "m", "./metrics", "path to metrics folder.")

}

func visualizeMetrics(cmd *cobra.Command, args []string) {
	metricsFolder, err := cmd.Flags().GetString(metricsFlag)
	if err != nil {
		panic(err)
	}
	metrics, err := loadMetricsFromFolder(metricsFolder)
	if err != nil {
		fmt.Printf("failed to load metrics from %s : %s", metricsFolder, err)
		return
	}
	fmt.Println(len(metrics))

}

func loadMetricsFromFolder(folder string) (results []*visMetric, err error) {
	items, err := ioutil.ReadDir(folder)
	if err != nil {
		return results, err
	}
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		m, err := loadMetricsFromFile(path.Join(folder, item.Name()))
		runUUID := strings.Split(item.Name(), ".json")[0]

		for _, mi := range m {
			mi.RunUUID = runUUID
		}
		if err != nil {
			return nil, err
		}
		results = append(results, m...)
	}
	return results, err
}

func loadMetricsFromFile(filePath string) (m []*visMetric, err error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return m, err
	}
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			fmt.Printf("failed to close json file %s", err)
		}
	}()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(byteValue, &m)
	return m, err
}
