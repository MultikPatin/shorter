package psql

const (
	createLinksTable = `
		CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		origin VARCHAR(255),
		short VARCHAR(255)
		);`
	addShortLink = `
		INSERT INTO events (short, origin) 
		VALUES ($1, $2)
		RETURNING id;`
	getShortLink = `
		SELECT origin 
		FROM events 
		WHERE short = $1;`
	getOrigin = `
		SELECT short 
		FROM events 
		WHERE origin = $1;`
)
