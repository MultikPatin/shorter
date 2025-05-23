package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"main/internal/constants"
	"main/internal/models"
	"main/internal/services"
)

// LinksRepository manages CRUD operations for links stored in a PostgreSQL database.
type LinksRepository struct {
	db *PostgresDB // Reference to the PostgreSQL database handler.
}

// NewLinksRepository constructs a new LinksRepository instance connected to a specific PostgresDB.
func NewLinksRepository(db *PostgresDB) *LinksRepository {
	return &LinksRepository{
		db: db,
	}
}

// Add inserts a new link record into the database, handling potential conflicts.
func (r *LinksRepository) Add(ctx context.Context, addedLink models.AddedLink) (string, error) {
	userID := ctx.Value(constants.UserIDKey).(int64)

	_, err := r.db.Connection.ExecContext(ctx, addShortLink, addedLink.Short, addedLink.Origin, userID)
	if err == nil {
		return addedLink.Short, nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || !pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		return "", err
	}

	var shortLink string
	err = r.db.Connection.QueryRowContext(ctx, getOrigin, addedLink.Origin).Scan(&shortLink)
	if err != nil {
		return "", err
	}

	return shortLink, services.ErrConflict
}

// AddBatch bulk-adds multiple links atomically using transactions.
func (r *LinksRepository) AddBatch(ctx context.Context, addedLinks []models.AddedLink) ([]models.Result, error) {
	userID := ctx.Value(constants.UserIDKey).(int64)

	tx, err := r.db.Connection.Begin()
	if err != nil {
		return nil, err
	}

	var results []models.Result

	for _, link := range addedLinks {
		_, err := r.db.Connection.ExecContext(ctx, addShortLink, link.Short, link.Origin, userID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		result := models.Result{
			CorrelationID: link.CorrelationID,
			Result:        link.Short,
		}
		results = append(results, result)
	}
	tx.Commit()
	return results, nil
}

// Get retrieves the original URL associated with a given short link.
func (r *LinksRepository) Get(ctx context.Context, short string) (string, error) {
	var originalLink string
	var isDeleted bool

	err := r.db.Connection.QueryRowContext(ctx, getShortLink, short).Scan(&originalLink, &isDeleted)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("link with short %s not found", short)
	} else if err != nil {
		return "", err
	}
	if isDeleted {
		return "", services.ErrDeletedLink
	}
	return originalLink, nil
}
