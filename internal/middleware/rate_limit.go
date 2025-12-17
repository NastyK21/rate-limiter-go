package middleware

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NastyK21/rate-limiter-go/internal/config"
	"github.com/NastyK21/rate-limiter-go/internal/limiter"
	"github.com/NastyK21/rate-limiter-go/internal/metrics"
	"github.com/NastyK21/rate-limiter-go/pkg/response"
)

type RateLimitConfig struct {
	Capacity        float64
	RefillRate      float64
	FailureStrategy config.FailureStrategy
	LocalLimiter    *limiter.LocalLimiter

	// Phase 6: per-user limits
	UserCapacity   float64
	UserRefillRate float64
}

func extractUserID(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", false
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return "", false
	}

	userID := strings.TrimPrefix(auth, prefix)
	if userID == "" {
		return "", false
	}

	return userID, true
}

func RateLimit(
	rl *limiter.RateLimiter,
	cfg RateLimitConfig,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			// 1️⃣ Extract client IP
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				response.WriteError(w, http.StatusBadRequest, "invalid client address")
				return
			}

			// 2️⃣ Timeout to protect API latency
			ctx, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
			defer cancel()

			var (
				key        string
				capacity   float64
				refillRate float64
				identity   string
			)

			// 3️⃣ Decide identity (user vs ip)
			if userID, ok := extractUserID(r); ok {
				key = "user:" + userID
				capacity = cfg.UserCapacity
				refillRate = cfg.UserRefillRate
				identity = "user"
			} else {
				key = "ip:" + normalizeIP(ip)
				capacity = cfg.Capacity
				refillRate = cfg.RefillRate
				identity = "ip"
			}

			// 4️⃣ Rate limit check (Redis)
			allowed, remaining, err := rl.Allow(ctx, key, capacity, refillRate)

			if err != nil {
				metrics.RateLimitErrors.Inc()

				if cfg.FailureStrategy == config.FailOpen {
					metrics.RateLimitDegraded.Inc()

					allowed := cfg.LocalLimiter.Allow(key)
					w.Header().Set("X-RateLimit-Degraded", "true")

					if !allowed {
						metrics.RateLimitBlocked.WithLabelValues(identity).Inc()
						response.WriteError(w, http.StatusTooManyRequests, "rate limit exceeded (degraded)")
						return
					}

					metrics.RateLimitAllowed.WithLabelValues(identity).Inc()
					next.ServeHTTP(w, r)
					return
				}

				metrics.RateLimitBlocked.WithLabelValues(identity).Inc()
				response.WriteError(w, http.StatusTooManyRequests, "rate limiter unavailable")
				return
			}

			// 5️⃣ Rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.FormatFloat(capacity, 'f', -1, 64))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatFloat(remaining, 'f', -1, 64))

			// 6️⃣ Block if limit exceeded
			if !allowed {
				metrics.RateLimitBlocked.WithLabelValues(identity).Inc()
				w.Header().Set("Retry-After", "1")
				response.WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}

			// 7️⃣ Allow request
			//log.Println("RATE LIMIT ALLOWED:", identity)
			metrics.RateLimitAllowed.WithLabelValues(identity).Inc()
			next.ServeHTTP(w, r)

		})
	}
}

func normalizeIP(ip string) string {
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}
