package app

type ShortenRequest struct {
	Url string `json:"url,omitempty"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}
