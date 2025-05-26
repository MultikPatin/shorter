package middleware

import (
	"main/internal/adapters"
	"net/http"
	"time"
)

// Structure to hold additional response metrics.
type responseData struct {
	status int // HTTP status code returned by the handler.
	size   int // Size of the response body written.
}

// Logging-aware response writer wrapper.
type loggingResponseWriter struct {
	http.ResponseWriter               // Delegate for actual response writing.
	responseData        *responseData // Tracks response metadata.
}

// Override Write to track the amount of data written.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader Override WriteHeader to capture the final status code.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// AccessLogger logs essential request and response metrics for every handled request.
func AccessLogger(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		logger := adapters.GetLogger()

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
