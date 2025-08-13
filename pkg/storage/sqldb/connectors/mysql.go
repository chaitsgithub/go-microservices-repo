package sqldb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

type mysqlconnector struct {
	cfg *DBConfig
}

func NewMySQLConnector(cfg *DBConfig) Connector {
	return &mysqlconnector{cfg: cfg}
}

func (m *mysqlconnector) Connect() (*sql.DB, error) {
	dsn := m.cfg.dsn(DB_MYSQL)
	db, err := sql.Open(DB_MYSQL, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySql DB connection. Error: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping MySql DB. Error: %w", err)
	}
	return db, nil
}
