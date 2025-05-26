package models

// AddedLink captures the result of a successful link addition operation.
type AddedLink struct {
	CorrelationID string // Identifier correlating with the originating request.
	Short         string // Generated short URL.
	Origin        string // Original long URL.
}

// OriginLink represents a link submission for shortening.
type OriginLink struct {
	CorrelationID string // Identifier for tracking purposes.
	URL           string // Long URL to be shortened.
}

// Result summarizes the outcome of a link shortening attempt.
type Result struct {
	CorrelationID string // Associated correlation ID.
	Result        string // Final short URL or error message.
}

// UserLinks pairs a short URL with its corresponding original URL.
type UserLinks struct {
	Shorten  string // Shortened URL.
	Original string // Original URL.
}
