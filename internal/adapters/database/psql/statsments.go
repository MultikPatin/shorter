package psql

const (
	createLinksTable = "CREATE TABLE IF NOT EXISTS events (id SERIAL PRIMARY KEY,original VARCHAR(255),short VARCHAR(255));"
)
