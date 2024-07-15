package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingWriter struct {
	http.ResponseWriter
	code int
}

func newLoggingWriter(w http.ResponseWriter) *loggingWriter {
	return &loggingWriter{
		ResponseWriter: w,
		code:           http.StatusInternalServerError,
	}
}

func (lw *loggingWriter) WriteHeader(code int) {
	lw.code = code
	lw.ResponseWriter.WriteHeader(code)
}

func UseAccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rlw := newLoggingWriter(w)

		next.ServeHTTP(rlw, req)

		slog.InfoContext(req.Context(), "access log",
			slog.String("URI", req.RequestURI),
			slog.String("Method", req.Method),
			slog.Int("Status code", rlw.code),
			slog.Int64("Elapsed time (ms)", time.Since(start).Milliseconds()),
		)
	})
}
