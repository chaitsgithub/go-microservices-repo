package sqldb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type postgressqlconnector struct {
	cfg *DBConfig
}

func NewPostgreSQLConnector(cfg *DBConfig) Connector {
	return &postgressqlconnector{cfg: cfg}
}

func (m *postgressqlconnector) Connect() (*sql.DB, error) {
	dsn := m.cfg.dsn(DB_POSTGRES)
	db, err := sql.Open(DB_POSTGRES, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySql DB connection. Error: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping MySql DB. Error: %w", err)
	}
	return db, nil
}
