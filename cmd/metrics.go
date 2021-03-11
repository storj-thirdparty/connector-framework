package cmd

import (
	"runtime"
	"time"
)

type Metric struct {
	function    string
	startMemory *runtime.MemStats
	endMemory   *runtime.MemStats
	startTime   int64
	endTime     int64
}

func (m *Metric) start() {
	m.startTime = time.Now().Unix()
	runtime.ReadMemStats(m.startMemory)
}

func (m *Metric) end() {
	m.endTime = time.Now().Unix()
	runtime.ReadMemStats(m.endMemory)
}
