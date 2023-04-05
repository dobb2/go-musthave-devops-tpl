package client

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dobb2/go-musthave-devops-tpl/internal/config"
	"github.com/dobb2/go-musthave-devops-tpl/internal/crypto"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/go-resty/resty/v2"
)

func SendMetric(metric metrics.Metrics, cfg config.AgentConfig) {
	client := resty.New().
		SetBaseURL("http://" + cfg.Address).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	out, err := json.Marshal(metric)
	if err != nil {
		log.Println(err)
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(out).
		Post("/update/")

	if err != nil {
		log.Println(err)
	}

	log.Println(resp.StatusCode())
}

func PutMetric(m *cache.Metrics, cfg config.AgentConfig) {
	for _, Metric := range m.Metrics { // Порядок не определен

		switch Metric.MType {
		case "counter":
			Metric.Hash = crypto.Hash(fmt.Sprintf("%s:counter:%d", Metric.ID, *Metric.Delta), cfg.Key)
		case "gauge":
			Metric.Hash = crypto.Hash(fmt.Sprintf("%s:gauge:%f", Metric.ID, *Metric.Value), cfg.Key)
		default:
			log.Println("invalid type metric for create hash")
		}

		SendMetric(Metric, cfg)
	}
}
