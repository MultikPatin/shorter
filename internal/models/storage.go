package models

type Event struct {
	ID     int    `json:"uuid"`
	Origin string `json:"original_url"`
	Short  string `json:"short_url"`
}
