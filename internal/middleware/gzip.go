package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	json = "application/json"
	html = "text/html"
)

var availableContentTypes = []string{json, html}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GZipper(h http.Handler) http.Handler {
	zipFn := func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentEncoding := r.Header.Get("Content-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		zipped := false

		requestContentType := r.Header.Get("Content-Type")
		for _, ContentType := range availableContentTypes {
			if ContentType == requestContentType {
				zipped = true
				break
			}
		}

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
