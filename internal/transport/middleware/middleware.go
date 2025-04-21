package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	statuscode int
}

func (w *wrappedWriter) WriteHeader(statuscode int) {
	w.ResponseWriter.WriteHeader(statuscode)
	w.statuscode = statuscode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statuscode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)
		slog.Info("Request made", 
			slog.Int("status_code", wrapped.statuscode),
			slog.String("method" , r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration",time.Since(start)),
			slog.String("client_ip",r.RemoteAddr),
		)
	})
}