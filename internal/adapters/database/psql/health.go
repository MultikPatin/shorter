package psql

type HealthRepository struct {
	db *PostgresDB
}

func NewHealthRepository(db *PostgresDB) *HealthRepository {
	return &HealthRepository{
		db: db,
	}
}

func (r *HealthRepository) Ping() error {
	err := r.db.Connection.Ping()
	return err
}
