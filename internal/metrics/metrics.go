package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RateLimitAllowed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_allowed_total",
			Help: "Total number of allowed requests",
		},
		[]string{"identity"},
	)

	RateLimitBlocked = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_blocked_total",
			Help: "Total number of blocked requests",
		},
		[]string{"identity"},
	)

	RateLimitDegraded = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "rate_limit_degraded_total",
			Help: "Total number of requests served in degraded mode",
		},
	)

	RateLimitErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "rate_limit_errors_total",
			Help: "Total number of rate limiter errors",
		},
	)
)

func Register() {
	prometheus.MustRegister(
		RateLimitAllowed,
		RateLimitBlocked,
		RateLimitDegraded,
		RateLimitErrors,
	)
}
