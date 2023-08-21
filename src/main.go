package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

type EndpointCheck struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url"`
	Frequency string `yaml:"frequency"`
}

type Config struct {
	Checks []EndpointCheck `yaml:"checks"`
}

var (
	latencyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_latency_milliseconds",
			Help: "Latency of HTTP requests to endpoints",
		},
		[]string{"endpoint"},
	)
	statusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_up",
			Help: "Whether the HTTP endpoint is up (1) or down (0)",
		},
		[]string{"endpoint"},
	)
)

func main() {
	// Read YAML file
	yamlData, err := ioutil.ReadFile("endpoints.yaml")
	if err != nil {
		panic(err)
	}

	// Parse YAML
	var config Config
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		panic(err)
	}

	// Initialize Prometheus metrics
	prometheus.MustRegister(latencyGauge)
	prometheus.MustRegister(statusGauge)

	// Create a channel to signal when the program should stop
	stopChan := make(chan struct{})

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start goroutines for each endpoint check
	for _, check := range config.Checks {
		wg.Add(1)
		go func(check EndpointCheck) {
			defer wg.Done()
			checkEndpoint(&check, stopChan)
		}(check)
	}

	// Start HTTP server to expose Prometheus metrics and custom endpoints
	http.Handle("/metricsz", promhttp.Handler())
	http.HandleFunc("/healthz", healthCheckHandler)
	go http.ListenAndServe(":9090", nil)

	// Wait for interrupt signal to stop the application
	fmt.Println("Press Ctrl+C to stop...")
	wg.Wait()
	close(stopChan)
}

func checkEndpoint(endpoint *EndpointCheck, stopChan <-chan struct{}) {
	frequency := parseDurationOrDefault(endpoint.Frequency, 30*time.Second)
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			startTime := time.Now()
			resp, err := http.Get(endpoint.URL)
			latency := time.Since(startTime).Milliseconds()

			if err != nil {
				fmt.Printf("Error checking %s: %v\n", endpoint.Name, err)
				statusGauge.WithLabelValues(endpoint.Name).Set(0)
				latencyGauge.WithLabelValues(endpoint.Name).Set(0) // Set latency to 0 on error
			} else {
				resp.Body.Close()
				fmt.Printf("%s: Latency: %dms\n", endpoint.Name, latency)
				statusGauge.WithLabelValues(endpoint.Name).Set(1)
				latencyGauge.WithLabelValues(endpoint.Name).Set(float64(latency))
			}
		case <-stopChan:
			return
		}
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func parseDurationOrDefault(durationStr string, defaultValue time.Duration) time.Duration {
	if durationStr == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		panic(fmt.Errorf("invalid duration: %s", durationStr))
	}
	return duration
}
