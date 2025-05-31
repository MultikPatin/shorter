package models

// Event tracks the history of link transformations.
type Event struct {
	Origin string `json:"original_url"` // Original URL being tracked.
	Short  string `json:"short_url"`    // Shortened equivalent of the original URL.
	ID     int    `json:"uuid"`         // Unique identifier for the event.
}
