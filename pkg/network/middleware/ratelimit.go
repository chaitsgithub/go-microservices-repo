package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func WithRateLimiter(requests int, duration time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if requests == 0 {
			requests = 100
		}
		if duration == 0 {
			duration = time.Minute
		}
		return httprate.LimitByIP(requests, duration)(next)
	}
}
