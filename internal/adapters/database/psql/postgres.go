package psql

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"net/url"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(PostgresDNS *url.URL, logger *zap.SugaredLogger) (*PostgresDB, error) {
	host := PostgresDNS.Hostname()
	port := PostgresDNS.Port()
	user := PostgresDNS.User.Username()
	password, _ := PostgresDNS.User.Password()
	dbname := PostgresDNS.Path[1:]

	ps := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("pgx", ps)
	if err != nil {
		logger.Infow(
			"Create PostgresDB",
			"error", err.Error(),
		)
	}
	postgresDB := PostgresDB{
		db: db,
	}
	return &postgresDB, err
}

func (p *PostgresDB) Close() error {
	p.db.Close()
	return nil
}

func (p *PostgresDB) Ping() error {
	err := p.db.Ping()
	return err
}

func (p *PostgresDB) Add(id string, link string) (string, error) {
	return "", nil
}

func (p *PostgresDB) Get(id string) (string, error) {
	return "", nil
}
