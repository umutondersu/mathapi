package middleware

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/umutondersu/mathapi/internal/ratelimit"
	"golang.org/x/time/rate"
)

func LimitMiddleware(next http.Handler) http.Handler {
	rate := rate.Limit(5)
	burst := 1
	limiter := ratelimit.NewIPRateLimiter(rate, burst)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			logger.Error("Rate limit exceeded", slog.String("ip", r.RemoteAddr))
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
