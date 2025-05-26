package memory

// HealthRepository serves as a repository layer abstraction providing health checks functionality.
type HealthRepository struct {
	db *InMemoryDB // Reference to the underlying in-memory database instance.
}

// NewHealthRepository constructs a new HealthRepository instance tied to an InMemoryDB.
func NewHealthRepository(db *InMemoryDB) *HealthRepository {
	return &HealthRepository{
		db: db,
	}
}

// Ping simulates a basic health check method, always returning success without actual validation logic.
func (db *HealthRepository) Ping() error {
	return nil
}
