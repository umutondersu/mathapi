package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/umutondersu/mathapi/internal/ratelimit"
	"golang.org/x/time/rate"
)

type contextKey string

const loggerKey = contextKey("logger")

func LimitMiddleware(next http.Handler) http.Handler {
	rate := rate.Limit(1)
	burst := 5
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

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		ctx := context.WithValue(r.Context(), loggerKey, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetLogger(r *http.Request) *slog.Logger {
	return r.Context().Value(loggerKey).(*slog.Logger)
}

func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
