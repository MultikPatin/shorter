package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Available content types that support gzip encoding.
var availableContentTypes = map[string]bool{
	"json": true,
	"html": true,
}

// gzipWriter adapts ResponseWriter to work with compressed output streams.
type gzipWriter struct {
	http.ResponseWriter           // Wraps the original response writer.
	Writer              io.Writer // Compressed stream writer.
}

// Write overrides the Write method to delegate to the wrapped writer.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GZipper compresses outgoing HTTP responses if the client accepts gzip encoding.
func GZipper(h http.Handler) http.Handler {
	zipFn := func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentEncoding := r.Header.Get("Content-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		requestContentType := r.Header.Get("Content-Type")

		zipped := availableContentTypes[requestContentType]

		if sendsGzip {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer gz.Close()
		}
		if supportsGzip && zipped {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")
			h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		} else {
			h.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(zipFn)
}
