package models

// ShortenRequest represents a single link shortening request.
type ShortenRequest struct {
	URL string `json:"url,omitempty"` // Optional field for the URL to be shortened.
}

// ShortenResponse carries the result of a link shortening operation.
type ShortenResponse struct {
	Result string `json:"result"` // Resulting short URL.
}

// ShortensRequest encapsulates a batch link shortening request item.
type ShortensRequest struct {
	CorrelationID string `json:"correlation_id,omitempty"` // Unique correlation ID for traceability.
	URL           string `json:"original_url,omitempty"`   // Original URL to be shortened.
}

// ShortensResponse conveys the result of a batch link shortening operation.
type ShortensResponse struct {
	CorrelationID string `json:"correlation_id,omitempty"` // Corresponding correlation ID.
	Result        string `json:"short_url"`                // Generated short URL.
}

// UserLinksResponse represents a user-facing link summary with both short and original URLs.
type UserLinksResponse struct {
	Shorten  string `json:"short_url"`    // Shortened URL.
	Original string `json:"original_url"` // Original URL.
}
