package main

import (
	"github.com/caarlos0/env/v7"
	"github.com/dobb2/go-musthave-devops-tpl/internal/client"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"log"
	"time"
)

func main() {
	var cfg config.AgentConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	m := cache.Create()

	ticker := time.NewTicker(cfg.PollInterval)
	ticker2 := time.NewTicker(cfg.ReportInterval)

	for {
		select {
		case <-ticker.C:
			m.CollectMetrics()
		case <-ticker2.C:
			client.PutMetric(&m, cfg)
		}
	}
}
