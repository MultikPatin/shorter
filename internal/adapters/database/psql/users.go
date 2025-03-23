package psql

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"main/internal/constants"
	"main/internal/models"
	"main/internal/services"
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

	rows, err := r.db.Connection.QueryContext(ctx, getLinksByUser, userID)
	if err != nil {
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
	if len(links) == 0 {
		return nil, services.ErrNoLinksByUser
	}
	return links, nil
}

func (r *UsersRepository) DeleteLinks(ctx context.Context, shortLinks []string) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	//var links []models.UserLinks
	//userID := ctx.Value(constants.UserIDKey).(int64)
	//
	//rows, err := r.db.Connection.QueryContext(ctx, getLinksByUser, userID)
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var link models.UserLinks
	//	err := rows.Scan(&link.Shorten, &link.Original)
	//	if err != nil {
	//		return nil, err
	//	}
	//	links = append(links, link)
	//}
	//if err := rows.Err(); err != nil {
	//	return nil, err
	//}
	//if len(links) == 0 {
	//	return nil, services.ErrNoLinksByUser
	//}
	return
}
