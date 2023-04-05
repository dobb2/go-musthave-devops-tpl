package main

import (
	"log"
	"time"

	"github.com/dobb2/go-musthave-devops-tpl/internal/client"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
)

func main() {
	cfg := config.CreateAgentConfig()
	m := cache.Create()
	log.Println("key ", cfg.Key)

	tickerCollector := time.NewTicker(cfg.PollInterval)
	tickerSender := time.NewTicker(cfg.ReportInterval)

	for {
		select {
		case <-tickerCollector.C:
			m.CollectMetrics()
		case <-tickerSender.C:
			client.PutMetric(&m, cfg)
		}
	}
}
