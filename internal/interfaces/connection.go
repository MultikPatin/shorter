package interfaces

// DB describes a minimal set of behaviors expected from a database driver or client library.
type DB interface {
	Close() error // Closes the database connection safely.
	Ping() error  // Tests the database connection by sending a lightweight request.
}
