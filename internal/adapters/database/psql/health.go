package psql

// HealthRepository abstracts the process of performing health checks against a PostgreSQL database.
type HealthRepository struct {
	db *PostgresDB // Reference to the PostgreSQL database wrapper.
}

// NewHealthRepository instantiates a new HealthRepository bound to a particular PostgresDB instance.
func NewHealthRepository(db *PostgresDB) *HealthRepository {
	return &HealthRepository{
		db: db,
	}
}

// Ping executes a simple health check by attempting to connect to the database.
func (r *HealthRepository) Ping() error {
	err := r.db.Ping()
	return err
}
