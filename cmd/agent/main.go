package main

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/client"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"time"
)

func main() {
	m := cache.Create()
	ticker := time.NewTicker(config.PollInterval * time.Second)
	ticker2 := time.NewTicker(config.ReportInterval * time.Second)

	for {
		select {
		case <-ticker.C:
			m.CollectMetrics()
		case <-ticker2.C:
			client.PutMetric(&m)
		}
	}
}
