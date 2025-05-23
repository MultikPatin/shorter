package psql

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"main/internal/constants"
	"main/internal/models"
	"main/internal/services"
	"time"
)

// UsersRepository manages user login, link retrieval, and deletion operations in a PostgreSQL database.
type UsersRepository struct {
	db *PostgresDB // Reference to the PostgreSQL database handler.
}

// NewUsersRepository constructs a new UsersRepository instance linked to a specific PostgresDB.
func NewUsersRepository(db *PostgresDB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

// Login registers a new user session and retrieves their assigned user ID.
func (r *UsersRepository) Login(ctx context.Context) (int64, error) {
	var userID int64

	err := r.db.Connection.QueryRowContext(ctx, addUser).Scan(&userID)
	if err != nil {
		return -1, fmt.Errorf("couldn't add user: %w", err)
	}
	return userID, nil
}

// GetLinks fetches all links created by a specific user.
func (r *UsersRepository) GetLinks(ctx context.Context) ([]models.UserLinks, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var links []models.UserLinks
	userID := ctx.Value(constants.UserIDKey).(int64)

	rows, err := r.db.Connection.QueryContext(ctx, getLinksByUser, userID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get the user's links: %w", err)
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

// DeleteLinks removes a list of short links belonging to a specific user.
func (r *UsersRepository) DeleteLinks(ctx context.Context, shortLinks []string) error {
	userID := ctx.Value(constants.UserIDKey).(int64)

	tx, err := r.db.Connection.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("couldn't start a transaction: %w", err)
	}

	for _, link := range shortLinks {
		_, err := tx.ExecContext(ctx, deleteLinksByUser, link, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("couldn't confirm the transaction: %w", err)
	}

	return nil
}
