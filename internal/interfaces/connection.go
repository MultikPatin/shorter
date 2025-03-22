package interfaces

type DB interface {
	Close() error
	Ping() error
}
