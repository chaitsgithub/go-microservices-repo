package repositories

import (
	"chaits.org/go-microservices-repo/pkg/general/logger"
	sqldb "chaits.org/go-microservices-repo/pkg/storage/sqldb/connectors"
)

// DBManager holds all table-specific repositories.
type DBManager struct {
	AppRepo AppRepository
}

// NewDBManager initializes the database connection and repositories.
func NewMySQLDBManager() (*DBManager, error) {
	dbconn, _ := sqldb.NewConnector(sqldb.DB_MYSQL)

	db, err := dbconn.Connect()
	if err != nil {
		logger.Logger.WithError(err).Error("error connecting to db")
	}

	// Initialize table-specific repositories
	appRepo := NewAppRepository(db)

	return &DBManager{
		AppRepo: appRepo,
	}, nil
}
