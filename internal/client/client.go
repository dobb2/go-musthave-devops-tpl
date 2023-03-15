package client

import (
	"bytes"
	"fmt"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

func SendMetric(TypeMetrics, NameMetric, ValueMetric string) {
	path := path.Join("update", TypeMetrics, NameMetric, ValueMetric)
	endpoint, err := url.Parse("http://127.0.0.1:8080/path")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	endpoint.Path = path
	var data = []byte(strings.Join([]string{NameMetric, ValueMetric}, ": "))
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	request.Header.Add("Content-Type", "text/plain")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(response.StatusCode)
	log.Println(string(body))
}

func PutMetric(m *cache.Metrics) {
	for NameMetric, ValueMetric := range m.GaugeMetrics { // Порядок не определен
		SendMetric("gauge", NameMetric, fmt.Sprintf("%f", ValueMetric))
	}
	for NameMetric, ValueMetric := range m.CounterMetrics { // Порядок не определен
		SendMetric("counter", NameMetric, strconv.FormatInt(ValueMetric, 10))
	}
}
