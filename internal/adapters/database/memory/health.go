package memory

type HealthRepository struct {
	db *InMemoryDB
}

func NewHealthRepository(db *InMemoryDB) *HealthRepository {
	return &HealthRepository{
		db: db,
	}
}

func (db *HealthRepository) Ping() error {
	return nil
}
