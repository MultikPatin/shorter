package psql

const (
	createLinksTable = `
		CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		origin VARCHAR(255) NOT NULL,
		short VARCHAR(255) NOT NULL
		);`
	addShortLink = `
		INSERT INTO events (short, origin) 
		VALUES ($1, $2)`
	getShortLink = `
		SELECT origin 
		FROM events 
		WHERE short = $1;`
	getOrigin = `
		SELECT short 
		FROM events 
		WHERE origin = $1;`
)
