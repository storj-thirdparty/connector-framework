package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"time"
)

type Metric struct {
	Function   string `json:"Function"`
	StartHeap  uint64 `json:"StartHeap"`
	StartStack uint64 `json:"StartStack"`
	EndHeap    uint64 `json:"EndHeap"`
	EndStack   uint64 `json:"EndStack"`
	StartTime  int64  `json:"StartTime"`
	EndTime    int64  `json:"EndTime"`
}

func (m *Metric) start() {
	m.StartTime = time.Now().Unix()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.StartHeap = ms.HeapInuse
	m.StartStack = ms.StackInuse
}

func (m *Metric) end() {
	m.EndTime = time.Now().Unix()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.EndHeap = ms.HeapInuse
	m.EndStack = ms.StackInuse
}

func saveCollectedMetrics(metrics []*Metric) error {
	if len(metrics) == 0 {
		return nil
	}
	p := "metrics"
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = os.Mkdir(p, 0700)
		if err != nil {
			return err
		}
	}

	byteArr, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	metricsPath := path.Join(p, fmt.Sprintf("%s.json", uuid.New().String()))
	err = ioutil.WriteFile(metricsPath, byteArr, 0644)

	fmt.Printf("metrcis saved to %s", metricsPath)
	return err

}
