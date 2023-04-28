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

func SendMetric(metrics []metrics.Metrics, cfg config.AgentConfig) {
	client := resty.New().
		SetBaseURL("http://" + cfg.Address).
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	out, err := json.Marshal(metrics)
	if err != nil {
		log.Println(err)
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(out).
		Post("/updates/")

	if err != nil {
		log.Println(err)
	}

	log.Println(resp.StatusCode())
}

func PutMetric(m *cache.Metrics, cfg config.AgentConfig) {
	metrics := make([]metrics.Metrics, 0, len(m.Metrics))
	for _, Metric := range m.Metrics { // Порядок не определен
		switch Metric.MType {
		case "counter":
			Metric.Hash = crypto.Hash(fmt.Sprintf("%s:counter:%d", Metric.ID, *Metric.Delta), cfg.Key)
		case "gauge":
			Metric.Hash = crypto.Hash(fmt.Sprintf("%s:gauge:%f", Metric.ID, *Metric.Value), cfg.Key)
		default:
			log.Println("invalid type metric for create hash")
		}
		metrics = append(metrics, Metric)
	}
	SendMetric(metrics, cfg)
}
