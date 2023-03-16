package client

import (
	"fmt"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"github.com/go-resty/resty/v2"
	"log"
	"strconv"
	"time"
)

func SendMetric(TypeMetric, NameMetric, ValueMetric string) {
	client := resty.New().
		SetBaseURL("http://127.0.0.1:8080/").
		SetRetryCount(2).
		SetRetryWaitTime(1 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"typeMetric":  TypeMetric,
			"nameMetric":  NameMetric,
			"valueMetric": ValueMetric,
		}).Post("/update/{typeMetric}/{nameMetric}/{valueMetric}")

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(resp.StatusCode())
}

func PutMetric(m *cache.Metrics) {
	for NameMetric, ValueMetric := range m.GaugeMetrics { // Порядок не определен
		SendMetric("gauge", NameMetric, strconv.FormatFloat(ValueMetric, 'f', -1, 64))
	}
	for NameMetric, ValueMetric := range m.CounterMetrics { // Порядок не определен
		SendMetric("counter", NameMetric, strconv.FormatInt(ValueMetric, 10))
	}
}
