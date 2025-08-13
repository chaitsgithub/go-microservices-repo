package sqldb

import "fmt"

var mysqlConfig = DBConfig{
	DBDriver:   DB_MYSQL,
	DBHost:     "localhost",
	DBPort:     "3306",
	DBUser:     "mysqluser",
	DBPassword: "password",
	DBName:     "microservicesdb",
}

var postgresConfig = DBConfig{
	DBDriver:   DB_POSTGRES,
	DBHost:     "localhost",
	DBPort:     "5432",
	DBUser:     "postgresuser",
	DBPassword: "password",
	DBName:     "microservicesdb",
}

func NewConnector(DBDriver string) (Connector, error) {
	switch DBDriver {
	case DB_MYSQL:
		return NewMySQLConnector(&mysqlConfig), nil
	case DB_POSTGRES:
		return NewPostgreSQLConnector(&postgresConfig), nil
	default:
		return nil, fmt.Errorf("unsupported db driver : %s", DBDriver)
	}
}
