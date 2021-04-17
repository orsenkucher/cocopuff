package log

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey int

const (
	logCtx ctxKey = iota
)

func Middleware(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := r.Context()
			reqID := middleware.GetReqID(ctx)
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}

			var fields []zapcore.Field
			if reqID != "" {
				fields = append(fields, zap.String("reqID", reqID))
			}

			fields = append(fields,
				zap.String("path", r.URL.Path),
				zap.String("from", r.RemoteAddr),
				zap.String("proto", r.Proto),
				zap.String("uri", r.RequestURI),
				zap.String("host", r.Host),
				zap.String("scheme", scheme),
				zap.String("method", r.Method),
			)
			l := l.With(fields...)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				latency := time.Since(start)
				fields := []zapcore.Field{
					zap.Duration("took", latency),
					zap.Int64("ns", latency.Nanoseconds()),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
				}

				l.Info("request completed", fields...)
			}()

			ctx = context.WithValue(ctx, logCtx, l)
			r = r.WithContext(ctx)
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func For(ctx context.Context) (*zap.Logger, bool) {
	log, ok := ctx.Value(logCtx).(*zap.Logger)
	return log, ok
}
