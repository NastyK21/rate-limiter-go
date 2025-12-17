package middleware

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/NastyK21/rate-limiter-go/internal/limiter"
	"github.com/NastyK21/rate-limiter-go/pkg/response"
)

type RateLimitConfig struct {
	Capacity   float64
	RefillRate float64
}

func RateLimit(
	rl *limiter.RateLimiter,
	cfg RateLimitConfig,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 1️⃣ Extract client IP
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				response.WriteError(w, http.StatusBadRequest, "invalid client address")
				return
			}

			// 2️⃣ Timeout to protect API latency
			ctx, cancel := context.WithTimeout(r.Context(), 100*time.Millisecond)
			defer cancel()

			key := "ip:" + ip

			// 3️⃣ Rate limit check
			allowed, remaining, err := rl.Allow(
				ctx,
				key,
				cfg.Capacity,
				cfg.RefillRate,
			)

			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, "rate limiter error")
				return
			}

			// 4️⃣ Rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.FormatFloat(cfg.Capacity, 'f', -1, 64))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatFloat(remaining, 'f', -1, 64))

			// 5️⃣ Block if limit exceeded
			if !allowed {
				w.Header().Set("Retry-After", "1")
				response.WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}

			// 6️⃣ Allow request
			next.ServeHTTP(w, r)
		})
	}
}
