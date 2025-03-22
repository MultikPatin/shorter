package psql

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type UsersRepository struct {
	db *PostgresDB
}

func NewUsersRepository(db *PostgresDB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (r *UsersRepository) Login(ctx context.Context) (int64, error) {
	var userID int64

	err := r.db.Connection.QueryRowContext(ctx, addUser).Scan(&userID)
	if err != nil {
		return -1, err
	}
	return userID, nil
}
