package sqldb

import (
	"database/sql"
	"fmt"
)

const (
	DB_MYSQL    = "mysql"
	DB_POSTGRES = "postgres"
)

type Connector interface {
	Connect() (*sql.DB, error)
}

type DBConfig struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func (c *DBConfig) dsn(driverName string) string {
	switch driverName {
	case DB_MYSQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	case DB_POSTGRES:
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
	default:
		return ""
	}
}
