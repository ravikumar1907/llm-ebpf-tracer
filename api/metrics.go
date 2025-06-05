package api

import (
	"net/http"
	"strings"

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

	mlockCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "llm_mlock_calls_total",
			Help: "Total mlock calls traced by eBPF",
		},
	)

	openatCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "llm_openat_gpu_calls_total",
			Help: "Total openat calls on GPU devices",
		},
	)

	callsByComm = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "llm_syscall_by_process",
			Help: "Syscall type counts by comm",
		},
		[]string{"syscall", "comm"},
	)
)

func init() {
	prometheus.MustRegister(mmapCounter)
	prometheus.MustRegister(mlockCounter)
	prometheus.MustRegister(openatCounter)
	prometheus.MustRegister(callsByComm)

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)
}

func Increment(syscall string, commRaw [16]byte) {
	comm := strings.TrimRight(string(commRaw[:]), "\x00")
	switch syscall {
	case "mmap":
		mmapCounter.Inc()
	case "mlock":
		mlockCounter.Inc()
	case "openat":
		openatCounter.Inc()
	}
	callsByComm.WithLabelValues(syscall, comm).Inc()
}
