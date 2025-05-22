package psql

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"time"
)

// PostgresDB encapsulates a SQL database connection for interacting with a PostgreSQL backend.
type PostgresDB struct {
	Connection *sql.DB // The active database connection.
}

// Close terminates the database connection cleanly.
func (p *PostgresDB) Close() error {
	err := p.Connection.Close()
	if err != nil {
		return err
	}
	return nil
}

// Ping verifies connectivity to the database by issuing a ping request.
func (p *PostgresDB) Ping() error {
	err := p.Connection.Ping()
	return err
}

// NewPostgresDB establishes a new PostgreSQL database connection using provided credentials.
func NewPostgresDB(PostgresDNS *url.URL, logger *zap.SugaredLogger) (*PostgresDB, error) {
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
		logger.Infow("Error creating PostgreSQL connection", "error", err.Error())
		return nil, err
	}

	err = migrate(ctx, conn)
	if err != nil {
		logger.Infow("Error during table creation", "error", err.Error())
		return nil, err
	}

	postgresDB := PostgresDB{
		Connection: conn,
	}
	return &postgresDB, nil
}

// migrate applies database schema migrations using the provided context and connection.
func migrate(ctx context.Context, conn *sql.DB) error {
	_, err := conn.ExecContext(ctx, createTables)
	if err != nil {
		return err
	}
	return nil
}
