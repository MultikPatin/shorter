package psql

const (
	createLinksTable = `
		CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		origin VARCHAR(255),
		short VARCHAR(255)
		);`
	addShortLink = `
		INSERT INTO links (short, origin) 
		VALUES ($1, $2) 
		RETURNING id;`
	getShortLink = `
		SELECT * 
		FROM links 
		WHERE short = $1;`
)
