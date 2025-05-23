package constants

import "os"

// Storage file parameters when using the 'In Memory' database.
const (
	// DefaultFilePermissions specifies default Unix-style file permission bits (rw-rw-rw-) for new files.
	DefaultFilePermissions = 0666

	// DefaultProducerFileFlags are the default OS-level flags for opening/writing data by producers.
	// Combination of O_RDWR (read/write mode), O_CREATE (create if not exists), and O_APPEND (append at end).
	DefaultProducerFileFlags = os.O_RDWR | os.O_CREATE | os.O_APPEND

	// DefaultConsumerFileFlags are the default OS-level flags for reading data by consumers.
	// Combination of O_RDONLY (read-only mode) and O_CREATE (create if not exists).
	DefaultConsumerFileFlags = os.O_RDONLY | os.O_CREATE
)
