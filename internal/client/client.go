package client

import (
	"encoding/json"
	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/go-resty/resty/v2"
	"log"
	"time"
)

func SendMetric(metric metrics.Metrics, cfg config.EnvConfig) {
	client := resty.New().
		SetBaseURL("http://" + cfg.Address).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	out, err := json.Marshal(metric)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(out).
		Post("/update/")

	if err != nil {
		log.Panic(err)
	}

	log.Println(resp.StatusCode())
}

func PutMetric(m *cache.Metrics, cfg config.EnvConfig) {
	for _, Metric := range m.Metrics { // Порядок не определен
		SendMetric(Metric, cfg)
	}
}
