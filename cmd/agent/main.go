package main

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/metrics"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"time"
)

func main() {
	m := storage.NewMetrics()
	ticker := time.NewTicker(config.PollInterval * time.Second)
	ticker2 := time.NewTicker(config.ReportInterval * time.Second)

	for {
		select {
		case <-ticker.C:
			metrics.UpdateMetrics(m)
		case <-ticker2.C:
			metrics.PutMetric(m)
		}
	}
}
