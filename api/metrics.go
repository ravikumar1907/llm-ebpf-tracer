package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mmapCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "llm_mmap_calls_total",
			Help: "Total mmap calls traced by eBPF",
		},
	)
)

func Init() {
	prometheus.MustRegister(mmapCounter)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)
}

func IncrementMmap() {
	mmapCounter.Inc()
}
