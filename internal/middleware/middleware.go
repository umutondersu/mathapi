package middleware

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/umutondersu/mathapi/internal/ratelimit"
	"golang.org/x/time/rate"
)

type contextKey string

type Middleware func(http.Handler) http.Handler

const loggerKey = contextKey("logger")

func Ratelimit(next http.Handler) http.Handler {
	rate := rate.Limit(1)
	burst := 5
	limiter := ratelimit.NewIPRateLimiter(rate, burst)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := GetLogger(r)
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			logger.Error("Rate limit exceeded", slog.String("ip", r.RemoteAddr))
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		ctx := context.WithValue(r.Context(), loggerKey, logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetLogger(r *http.Request) *slog.Logger {
	logger, ok := r.Context().Value(loggerKey).(*slog.Logger)
	if !ok {
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return logger
}

func CreateStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			x := middlewares[i]
			next = x(next)
		}
		return next
	}
}
