package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type visMetric struct {
	Metric
	RunUUID string
}

type funcAvg struct {
	FuncName string
	Avg      float64
}

type visPage struct {
	HeapStartAvg  []*funcAvg
	HeapEndAvg    []*funcAvg
	HeapDeltaAvg  []*funcAvg
	StackStartAvg []*funcAvg
	EndStackAvg   []*funcAvg
	StackDeltaAvg []*funcAvg
	TimeSpendAvg  []*funcAvg
}

func visPageFromMetrics(metrics []*visMetric) *visPage {
	funcCounter := make(map[string]int64)
	hs := make(map[string]float64)
	he := make(map[string]float64)
	hd := make(map[string]float64)
	ss := make(map[string]float64)
	se := make(map[string]float64)
	sd := make(map[string]float64)
	ts := make(map[string]float64)

	for _, m := range metrics {
		funcCounter[m.Function] = funcCounter[m.Function] + 1
		hs[m.Function] = hs[m.Function] + float64(m.StartHeap)
		he[m.Function] = he[m.Function] + float64(m.EndHeap)
		hd[m.Function] = hd[m.Function] + float64(m.EndHeap-m.StartHeap)
		ss[m.Function] = ss[m.Function] + float64(m.StartStack)
		se[m.Function] = se[m.Function] + float64(m.EndStack)
		sd[m.Function] = sd[m.Function] + float64(m.EndStack-m.StartStack)
		ts[m.Function] = ts[m.Function] + float64(m.EndTime-m.StartTime)
	}
	var p visPage
	for funcName, count := range funcCounter {
		p.HeapStartAvg = append(p.HeapStartAvg, &funcAvg{
			FuncName: funcName,
			Avg:      hs[funcName] / float64(count),
		})
		p.HeapEndAvg = append(p.HeapEndAvg, &funcAvg{
			FuncName: funcName,
			Avg:      he[funcName] / float64(count),
		})
		p.HeapDeltaAvg = append(p.HeapDeltaAvg, &funcAvg{
			FuncName: funcName,
			Avg:      hd[funcName] / float64(count),
		})
		p.StackStartAvg = append(p.StackStartAvg, &funcAvg{
			FuncName: funcName,
			Avg:      ss[funcName] / float64(count),
		})
		p.EndStackAvg = append(p.EndStackAvg, &funcAvg{
			FuncName: funcName,
			Avg:      se[funcName] / float64(count),
		})
		p.StackDeltaAvg = append(p.StackDeltaAvg, &funcAvg{
			FuncName: funcName,
			Avg:      sd[funcName] / float64(count),
		})
		p.TimeSpendAvg = append(p.TimeSpendAvg, &funcAvg{
			FuncName: funcName,
			Avg:      ts[funcName] / float64(count),
		})

	}

	return &p
}

const metricsFlag = "metrics"

var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Command to visualize performance metrics.",
	Long:  `Command to visualize locally stored performance metrics.`,
	Run:   visualizeMetrics,
}

const templateSrc = `
<div id="chart"></div>
<script src="https://cdn.jsdelivr.net/npm/apexcharts"></script>
<script>
var point = {{.}};
alert(point);
</script>
<script>
var options = {
          series: [{
          name: 'series1',
          data: [31, 40, 28, 51, 42, 109, 100]
        }, {
          name: 'series2',
          data: [11, 32, 45, 32, 34, 52, 41]
        }],
          chart: {
          height: 350,
          type: 'area'
        },
        dataLabels: {
          enabled: false
        },
        stroke: {
          curve: 'smooth'
        },
        xaxis: {
          type: 'datetime',
          categories: ["2018-09-19T00:00:00.000Z", "2018-09-19T01:30:00.000Z", "2018-09-19T02:30:00.000Z", "2018-09-19T03:30:00.000Z", "2018-09-19T04:30:00.000Z", "2018-09-19T05:30:00.000Z", "2018-09-19T06:30:00.000Z"]
        },
        tooltip: {
          x: {
            format: 'dd/MM/yy HH:mm'
          },
        },
        };

        var chart = new ApexCharts(document.getElementById("chart"), options);
        chart.render();
</script>`

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

	http.HandleFunc("/", metricHandler(metrics))
	http.ListenAndServe(":8090", nil)

}

func metricHandler(metrics []*visMetric) func(http.ResponseWriter, *http.Request) {
	mp := visPageFromMetrics(metrics)
	pageJson, err := json.Marshal(mp)
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("").Parse(templateSrc))
	return func(w http.ResponseWriter, r *http.Request) {
		if err = t.Execute(w, string(pageJson)); err != nil {
			panic(err)
		}
	}
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
