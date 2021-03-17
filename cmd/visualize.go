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

type visPage struct {
	FuncNames     []string
	HeapStartAvg  []float64
	HeapEndAvg    []float64
	HeapDeltaAvg  []float64
	StackStartAvg []float64
	EndStackAvg   []float64
	StackDeltaAvg []float64
	TimeSpendAvg  []float64
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
		p.FuncNames = append(p.FuncNames, funcName)
		p.HeapStartAvg = append(p.HeapStartAvg, hs[funcName]/float64(count))
		p.HeapEndAvg = append(p.HeapEndAvg, he[funcName]/float64(count))
		p.HeapDeltaAvg = append(p.HeapDeltaAvg, hd[funcName]/float64(count))
		p.StackStartAvg = append(p.StackStartAvg, ss[funcName]/float64(count))
		p.EndStackAvg = append(p.EndStackAvg, se[funcName]/float64(count))
		p.StackDeltaAvg = append(p.StackDeltaAvg, sd[funcName]/float64(count))
		p.TimeSpendAvg = append(p.TimeSpendAvg, (ts[funcName]/float64(count))/float64(1000000))

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
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Connector Framework Metrics</title>
    <script src="https://cdn.jsdelivr.net/npm/apexcharts"></script>
</head>
<body>
<div>
    <h2>Time Spent</h2>
    <div id="chart-time"></div>
</div>
<div>
    <h2>Heap</h2>
    <div id="chart-heap"></div>
</div>
<div>
    <h2>Stack</h2>
    <div id="chart-stack"></div>
</div>

<script>
    const rawData = {{.}};
    const metrics = JSON.parse(rawData);
    console.log(metrics);

    const optionsTime = {
        series: [{
            name: 'Time Spent',
            data: metrics['TimeSpendAvg']
        },],
        chart: {
            type: 'bar',
            height: 350
        },
        plotOptions: {
            bar: {
                horizontal: false,
                columnWidth: '55%',
                endingShape: 'rounded'
            },
        },
        dataLabels: {
            enabled: false
        },
        stroke: {
            show: true,
            width: 2,
            colors: ['transparent']
        },
        xaxis: {
            categories: metrics['FuncNames'],
        },
        yaxis: {
            title: {
                text: 'ms'
            }
        },
        fill: {
            opacity: 1
        },
        tooltip: {
            y: {
                formatter: function (val) {
                    return val + " milliseconds"
                }
            }
        }
    };
    const chartTime = new ApexCharts(document.getElementById("chart-time"), optionsTime);
    chartTime.render();

    const optionsHeap = {
        series: [{
            name: 'Start Heap',
            data: metrics['HeapStartAvg']
        }, {
            name: 'End Heap',
            data: metrics['HeapEndAvg']
        }, {
            name: 'Delta Heap',
            data: metrics['HeapDeltaAvg']
        }],
        chart: {
            type: 'bar',
            height: 350
        },
        plotOptions: {
            bar: {
                horizontal: false,
                columnWidth: '55%',
                endingShape: 'rounded'
            },
        },
        dataLabels: {
            enabled: false
        },
        stroke: {
            show: true,
            width: 2,
            colors: ['transparent']
        },
        xaxis: {
            categories: metrics['FuncNames'],
        },
        yaxis: {
            title: {
                text: 'MB'
            }
        },
        fill: {
            opacity: 1
        },
        tooltip: {
            y: {
                formatter: function (val) {
                    return val + " MB"
                }
            }
        }
    };
    const chartHeap = new ApexCharts(document.getElementById("chart-heap"), optionsHeap);
    chartHeap.render();

    const optionsStack = {
        series: [{
            name: 'Start Stack',
            data: metrics['StackStartAvg']
        }, {
            name: 'End Stack',
            data: metrics['EndStackAvg']
        }, {
            name: 'Delta Stack',
            data: metrics['StackDeltaAvg']
        }],
        chart: {
            type: 'bar',
            height: 350
        },
        plotOptions: {
            bar: {
                horizontal: false,
                columnWidth: '55%',
                endingShape: 'rounded'
            },
        },
        dataLabels: {
            enabled: false
        },
        stroke: {
            show: true,
            width: 2,
            colors: ['transparent']
        },
        xaxis: {
            categories: metrics['FuncNames'],
        },
        yaxis: {
            title: {
                text: 'MB'
            }
        },
        fill: {
            opacity: 1
        },
        tooltip: {
            y: {
                formatter: function (val) {
                    return val + " MB"
                }
            }
        }
    };
    const chartStack = new ApexCharts(document.getElementById("chart-stack"), optionsStack);
    chartStack.render();

</script>
</body>
</html>
`

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
