package models

type ShortenRequest struct {
	URL string `json:"url,omitempty"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}
