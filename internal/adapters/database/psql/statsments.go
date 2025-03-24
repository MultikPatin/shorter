package psql

const (
	createTables = `
		CREATE TABLE IF NOT EXISTS users (
		    id SERIAL PRIMARY KEY
		);
		CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		origin VARCHAR(255) NOT NULL UNIQUE,
		short VARCHAR(255) NOT NULL,
		is_deleted BOOLEAN DEFAULT FALSE);
		CREATE INDEX IF NOT EXISTS origin_index ON events(origin);`
	// Links
	addShortLink = `
		INSERT INTO events (short, origin, user_id) 
		VALUES ($1, $2, $3)`
	getShortLink = `
		SELECT origin, is_deleted 
		FROM events 
		WHERE short = $1;`
	getOrigin = `
		SELECT short 
		FROM events 
		WHERE origin = $1;`
	// Users
	addUser = `
		INSERT INTO users DEFAULT VALUES RETURNING id;`
	getLinksByUser = `
		SELECT short, origin FROM events WHERE user_id = $1;`
	deleteLinksByUser = `
		UPDATE events SET is_deleted = true WHERE short = $1 AND user_id = $2;`
)
