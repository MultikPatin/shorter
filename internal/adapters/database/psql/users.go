package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"main/internal/constants"
	"main/internal/models"
	"time"
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

func (r *UsersRepository) GetLinks(ctx context.Context) ([]models.UserLinks, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var links []models.UserLinks
	userID := ctx.Value(constants.UserIDKey).(int64)
	fmt.Println(userID)

	rows, err := r.db.Connection.QueryContext(ctx, getLinksByUser, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("links by userID %d not found", userID)
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var link models.UserLinks
		err := rows.Scan(&link.Shorten, &link.Original)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return links, nil
}
