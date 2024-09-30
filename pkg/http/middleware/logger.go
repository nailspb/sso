package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"sso/pkg/helpers/slogHelper"
	"sso/pkg/http/routing"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logging(log *slog.Logger) routing.MiddlewareFunc {
	log = slogHelper.ConfigureForMiddleware(log, "Logging")
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			logger := slogHelper.AddRequestId(log, req.Context())
			ip := req.Header.Get("X-Real-IP")
			if ip == "" {
				ip = req.RemoteAddr
			}
			logger.Info("Request start",
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("query", req.URL.RawQuery),
				slog.String("ip", ip),
				slog.String("user-agent", req.UserAgent()),
			)
			recorder := &statusRecorder{
				ResponseWriter: w,
				Status:         200,
			}
			next.ServeHTTP(recorder, req)
			duration := time.Since(start)

			logger.Info("Request end",
				slog.Int("StatusCode", recorder.Status),
				slog.String("duration", fmt.Sprintf("%d us", duration.Microseconds())),
			)

		}
	}
}
