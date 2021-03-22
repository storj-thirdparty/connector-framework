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

const templateSrc = "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Connector Framework Metrics</title>\n    <!-- CSS only -->\n    <link href=\"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta2/dist/css/bootstrap.min.css\" rel=\"stylesheet\"\n          integrity=\"sha384-BmbxuPwQa2lc/FVzBcNJ7UAyJxM6wuqIj61tLrc4wSX0szH/Ev+nYRRuWlolflfl\" crossorigin=\"anonymous\">\n    <script src=\"https://cdn.jsdelivr.net/npm/apexcharts\"></script>\n</head>\n<body>\n<div class=\"container-fluid\">\n    <div >\n        <h2 class=\"d-flex justify-content-center\">Time Spent</h2>\n        <div id=\"chart-time\"></div>\n    </div>\n    <div>\n        <h2 class=\"d-flex justify-content-center\">Heap</h2>\n        <div id=\"chart-heap\"></div>\n    </div>\n    <div>\n        <h2 class=\"d-flex justify-content-center\">Stack</h2>\n        <div id=\"chart-stack\"></div>\n    </div>\n    <div class=\"mt-5\">\n        <table class=\"table table-striped\">\n            <thead class=\"thead thead-dark\">\n            <tr>\n                <th>Function Name</th>\n                <th>Time Spent, ms</th>\n                <th>Heap Start, MB</th>\n                <th>Heap End, MB</th>\n                <th>Heap Delta, MB</th>\n                <th>Stack Start, MB</th>\n                <th>Stack End, MB</th>\n                <th>Stack Delta, MB</th>\n            </tr>\n            </thead>\n            <tbody id=\"table-body\">\n\n            </tbody>\n        </table>\n    </div>\n</div>\n\n<script>\n    const rawData = \"{{.}}\";\n    const metrics = JSON.parse(rawData);\n    console.log(metrics)\n\n    const optionsTime = {\n        series: [{\n            name: 'Time Spent',\n            data: metrics['TimeSpendAvg']\n        },],\n        chart: {\n            type: 'bar',\n            height: 350\n        },\n        plotOptions: {\n            bar: {\n                horizontal: false,\n                columnWidth: '55%',\n                endingShape: 'rounded'\n            },\n        },\n        dataLabels: {\n            enabled: false\n        },\n        stroke: {\n            show: true,\n            width: 2,\n            colors: ['transparent']\n        },\n        xaxis: {\n            categories: metrics['FuncNames'],\n        },\n        yaxis: {\n            title: {\n                text: 'ms'\n            }\n        },\n        fill: {\n            opacity: 1\n        },\n        tooltip: {\n            y: {\n                formatter: function (val) {\n                    return val + \" milliseconds\"\n                }\n            }\n        }\n    };\n    const chartTime = new ApexCharts(document.getElementById(\"chart-time\"), optionsTime);\n    chartTime.render();\n\n    const optionsHeap = {\n        series: [{\n            name: 'Start Heap',\n            data: metrics['HeapStartAvg']\n        }, {\n            name: 'End Heap',\n            data: metrics['HeapEndAvg']\n        }, {\n            name: 'Delta Heap',\n            data: metrics['HeapDeltaAvg']\n        }],\n        chart: {\n            type: 'bar',\n            height: 350\n        },\n        plotOptions: {\n            bar: {\n                horizontal: false,\n                columnWidth: '55%',\n                endingShape: 'rounded'\n            },\n        },\n        dataLabels: {\n            enabled: false\n        },\n        stroke: {\n            show: true,\n            width: 2,\n            colors: ['transparent']\n        },\n        xaxis: {\n            categories: metrics['FuncNames'],\n        },\n        yaxis: {\n            title: {\n                text: 'MB'\n            }\n        },\n        fill: {\n            opacity: 1\n        },\n        tooltip: {\n            y: {\n                formatter: function (val) {\n                    return val + \" MB\"\n                }\n            }\n        }\n    };\n    const chartHeap = new ApexCharts(document.getElementById(\"chart-heap\"), optionsHeap);\n    chartHeap.render();\n\n    const optionsStack = {\n        series: [{\n            name: 'Start Stack',\n            data: metrics['StackStartAvg']\n        }, {\n            name: 'End Stack',\n            data: metrics['EndStackAvg']\n        }, {\n            name: 'Delta Stack',\n            data: metrics['StackDeltaAvg']\n        }],\n        chart: {\n            type: 'bar',\n            height: 350\n        },\n        plotOptions: {\n            bar: {\n                horizontal: false,\n                columnWidth: '55%',\n                endingShape: 'rounded'\n            },\n        },\n        dataLabels: {\n            enabled: false\n        },\n        stroke: {\n            show: true,\n            width: 2,\n            colors: ['transparent']\n        },\n        xaxis: {\n            categories: metrics['FuncNames'],\n        },\n        yaxis: {\n            title: {\n                text: 'MB'\n            }\n        },\n        fill: {\n            opacity: 1\n        },\n        tooltip: {\n            y: {\n                formatter: function (val) {\n                    return val + \" MB\"\n                }\n            }\n        }\n    };\n    const chartStack = new ApexCharts(document.getElementById(\"chart-stack\"), optionsStack);\n    chartStack.render();\n\n    let tableRowsHtml = '';\n    for (let i = 0; i < metrics.FuncNames.length; i++) {\n        tableRowsHtml += `\n           <tr>\n           <td>${metrics.FuncNames[i]}</td>\n           <td>${metrics.TimeSpendAvg[i]}</td>\n           <td>${metrics.HeapStartAvg[i]}</td>\n           <td>${metrics.HeapEndAvg[i]}</td>\n           <td>${metrics.HeapDeltaAvg[i]}</td>\n           <td>${metrics.StackStartAvg[i]}</td>\n           <td>${metrics.EndStackAvg[i]}</td>\n           <td>${metrics.StackDeltaAvg[i]}</td>\n           </tr>\n        `\n    }\n    const tableBody = document.getElementById(\"table-body\");\n    tableBody.innerHTML = tableRowsHtml;\n</script>\n</body>\n</html>\n\n\n\n"

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
