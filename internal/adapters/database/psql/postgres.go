package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"main/internal/models"
	"net/url"
	"time"
)

var ErrConflict = errors.New("data conflict")

type PostgresDB struct {
	conn *sql.DB
}

func NewPostgresRepository(PostgresDNS *url.URL, logger *zap.SugaredLogger) (*PostgresDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	host := PostgresDNS.Hostname()
	port := PostgresDNS.Port()
	user := PostgresDNS.User.Username()
	password, _ := PostgresDNS.User.Password()
	dbname := PostgresDNS.Path[1:]

	ps := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("pgx", ps)
	if err != nil {
		logger.Infow(
			"Create PostgresDB",
			"error", err.Error(),
		)
		return nil, err
	}

	err = migrate(ctx, conn)
	if err != nil {
		logger.Infow(
			"Create links table",
			"error", err.Error(),
		)
		return nil, err
	}
	postgresDB := PostgresDB{
		conn: conn,
	}
	return &postgresDB, err
}

func migrate(ctx context.Context, conn *sql.DB) error {
	_, err := conn.ExecContext(ctx, createLinksTable)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) Close() error {
	err := p.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) Ping() error {
	err := p.conn.Ping()
	return err
}

func (p *PostgresDB) Add(ctx context.Context, addedLink models.AddedLink) (string, error) {
	var shortLink string

	_, err := p.conn.ExecContext(ctx, addShortLink, addedLink.Short, addedLink.Origin)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err := p.conn.QueryRowContext(ctx, getOrigin, addedLink.Origin).Scan(&shortLink)
			if err != nil {
				return "", err
			}
			return shortLink, nil
		}
	}
	return addedLink.Short, ErrConflict
}

func (p *PostgresDB) AddBatch(ctx context.Context, addedLinks []models.AddedLink) ([]models.Result, error) {
	tx, err := p.conn.Begin()
	if err != nil {
		return nil, err
	}

	var results []models.Result

	for _, link := range addedLinks {
		_, err := p.conn.ExecContext(ctx, addShortLink, link.Short, link.Origin)
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

func (p *PostgresDB) Get(ctx context.Context, short string) (string, error) {
	var originalLink string
	err := p.conn.QueryRowContext(ctx, getShortLink, short).Scan(&originalLink)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("link with short %s not found", short)
	} else if err != nil {
		return "", err
	}
	return originalLink, nil
}
