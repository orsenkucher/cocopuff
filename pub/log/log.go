package log

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Middleware(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := middleware.GetReqID(r.Context())
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				latency := time.Since(start)
				fields := []zapcore.Field{
					zap.Duration("took", latency),
					zap.Int64("ns", latency.Nanoseconds()),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("path", r.URL.Path),
					zap.String("from", r.RemoteAddr),
					zap.String("proto", r.Proto),
					zap.String("uri", r.RequestURI),
					zap.String("host", r.Host),
					zap.String("scheme", scheme),
					zap.String("method", r.Method),
				}
				if reqID != "" {
					fields = append(fields, zap.String("reqID", reqID))
				}

				l.Info("request completed", fields...)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
