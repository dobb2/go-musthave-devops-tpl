package metrics

import (
	"bytes"
	"fmt"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func UpdateMetrics(m *storage.Metrics) {
	rand.Seed(time.Now().UnixNano())
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	// Misc memory stats
	m.GaugeMetrics["Alloc"] = float64(rtm.Alloc)
	m.GaugeMetrics["TotalAlloc"] = float64(rtm.TotalAlloc)
	m.GaugeMetrics["Sys"] = float64(rtm.Sys)
	m.GaugeMetrics["Mallocs"] = float64(rtm.Mallocs)
	m.GaugeMetrics["Frees"] = float64(rtm.Frees)
	m.GaugeMetrics["Alloc"] = float64(rtm.Alloc)
	m.GaugeMetrics["BuckHashSys"] = float64(rtm.BuckHashSys)
	m.GaugeMetrics["Frees"] = float64(rtm.Frees)
	m.GaugeMetrics["GCCPUFraction"] = float64(rtm.GCCPUFraction)
	m.GaugeMetrics["GCSys"] = float64(rtm.GCSys)
	m.GaugeMetrics["HeapAlloc"] = float64(rtm.HeapAlloc)
	m.GaugeMetrics["HeapIdle"] = float64(rtm.HeapIdle)
	m.GaugeMetrics["HeapInuse"] = float64(rtm.HeapInuse)
	m.GaugeMetrics["HeapObjects"] = float64(rtm.HeapObjects)
	m.GaugeMetrics["HeapReleased"] = float64(rtm.HeapReleased)
	m.GaugeMetrics["HeapSys"] = float64(rtm.HeapSys)
	m.GaugeMetrics["LastGC"] = float64(rtm.LastGC)
	m.GaugeMetrics["Lookups"] = float64(rtm.Lookups)
	m.GaugeMetrics["MCacheInuse"] = float64(rtm.MCacheInuse)
	m.GaugeMetrics["MCacheSys"] = float64(rtm.MCacheSys)
	m.GaugeMetrics["MSpanInuse"] = float64(rtm.MSpanInuse)
	m.GaugeMetrics["MSpanSys"] = float64(rtm.MSpanSys)
	m.GaugeMetrics["Mallocs"] = float64(rtm.Mallocs)
	m.GaugeMetrics["NextGC"] = float64(rtm.NextGC)
	m.GaugeMetrics["NumForcedGC"] = float64(rtm.NumForcedGC)
	m.GaugeMetrics["NumGC"] = float64(rtm.NumGC)
	m.GaugeMetrics["OtherSys"] = float64(rtm.OtherSys)
	m.GaugeMetrics["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	m.GaugeMetrics["StackInuse"] = float64(rtm.StackInuse)
	m.GaugeMetrics["StackSys"] = float64(rtm.StackSys)
	m.GaugeMetrics["Sys"] = float64(rtm.Sys)
	m.GaugeMetrics["TotalAlloc"] = float64(rtm.TotalAlloc)
	m.CounterMetrics["PollCount"]++
	m.GaugeMetrics["RandomValue"] = float64(rand.Float64())
}

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
	log.Println(string(body))
}

func PutMetric(m *storage.Metrics) {
	for NameMetric, ValueMetric := range m.GaugeMetrics { // Порядок не определен
		SendMetric("gauge", NameMetric, fmt.Sprintf("%f", ValueMetric))
	}
	for NameMetric, ValueMetric := range m.CounterMetrics { // Порядок не определен
		SendMetric("counter", NameMetric, strconv.FormatInt(ValueMetric, 10))
	}
}
