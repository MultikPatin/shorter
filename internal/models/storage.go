package models

type Event struct {
	ID       int    `json:"uuid"`
	Original string `json:"original_url"`
	Short    string `json:"short_url"`
}
