package models

type ShortenRequest struct {
	URL string `json:"url,omitempty"`
}
type ShortenResponse struct {
	Result string `json:"result"`
}

type ShortensRequest struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	URL           string `json:"original_url,omitempty"`
}

type ShortensResponse struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	Result        string `json:"short_url"`
}

type UserLinksResponse struct {
	Shorten  string `json:"short_url"`
	Original string `json:"original_url"`
}
