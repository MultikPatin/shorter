package constants

import "time"

// Content types and other common constants.
const (
	// ServerShutdownTime specifies the grace period for shutting down the server.
	ServerShutdownTime = 5 * time.Second

	// TextContentType is the MIME type for plain text content encoded in UTF-8.
	TextContentType = "text/plain; charset=utf-8"

	// JSONContentType is the MIME type for JSON-formatted data.
	JSONContentType = "application/json"

	// URLPrefix is the standard prefix for HTTP URLs.
	URLPrefix = "http://"

	// Delimiter is the forward slash character commonly used as a path separator in URLs.
	Delimiter = "/"
)
