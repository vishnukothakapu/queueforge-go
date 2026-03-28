package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	TotalJobs = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_total",
			Help: "Total number of jobs processed",
		},
	)

	FailedJobs = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_failed_total",
			Help: "Total number of failed jobs",
		},
	)

	RetriedJobs = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_retried_total",
			Help: "Total number of retried jobs",
		},
	)
)

func Init() {
	prometheus.MustRegister(TotalJobs)
	prometheus.MustRegister(FailedJobs)
	prometheus.MustRegister(RetriedJobs)
}
