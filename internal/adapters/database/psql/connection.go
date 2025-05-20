package psql

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"time"
)

type PostgresDB struct {
	Connection *sql.DB
}

func (p *PostgresDB) Close() error {
	err := p.Connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) Ping() error {
	err := p.Connection.Ping()
	return err
}

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
		logger.Infow(
			"Create Postgres Connection",
			"error", err.Error(),
		)
	}

	err = migrate(ctx, conn)
	if err != nil {
		logger.Infow(
			"Create tables",
			"error", err.Error(),
		)
	}
	postgresDB := PostgresDB{
		Connection: conn,
	}
	return &postgresDB, err
}

func migrate(ctx context.Context, conn *sql.DB) error {
	_, err := conn.ExecContext(ctx, createTables)
	if err != nil {
		return err
	}
	return nil
}
