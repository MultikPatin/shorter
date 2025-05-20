package constants

import "os"

const (
	DefaultFilePermissions   = 0666
	DefaultProducerFileFlags = os.O_RDWR | os.O_CREATE | os.O_APPEND
	DefaultConsumerFileFlags = os.O_RDONLY | os.O_CREATE
)
