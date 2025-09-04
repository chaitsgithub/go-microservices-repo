package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"chaits.org/go-microservices-repo/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
)

// AppRepository defines the interface for App-related database operations.
type AppRepository interface {
	CreateApp(ctx context.Context, a *models.App) (models.App, error)
	GetAllApps(ctx context.Context) ([]models.App, error)
	DeleteApp(ctx context.Context, id int) (sql.Result, error)
	ValidateAPIKey(ctx context.Context, appName, apiKey string) (string, bool, error)
}

// appRepository implements the AppRepository interface.
type appRepository struct {
	db *sql.DB
}

// NewAppRepository creates a new AppRepository.
func NewAppRepository(db *sql.DB) AppRepository {
	return &appRepository{db: db}
}

// CreateApp hashes the API key and inserts a new app into the database.
func (r *appRepository) CreateApp(ctx context.Context, a *models.App) (models.App, error) {
	_, tableSpan := otel.Tracer("db-tracer").Start(ctx, "db.insert.app")
	defer tableSpan.End()

	newApp := models.App{}
	// Hash the API key before storing it
	hashedAPIKey, err := bcrypt.GenerateFromPassword([]byte(a.APIKey), bcrypt.DefaultCost)
	if err != nil {
		return newApp, fmt.Errorf("failed to hash API key: %v", err)
	}

	// Insert into the database
	res, err := r.db.Exec("INSERT INTO apps (name, api_key_hash) VALUES (?, ?)", a.Name, hashedAPIKey)
	if err != nil {
		return newApp, fmt.Errorf("error inserting into db: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return newApp, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	newApp = models.App{
		ID:     int(id),
		Name:   a.Name,
		APIKey: a.APIKey,
	}

	return newApp, nil
}

// GetAllApps retrieves all registered apps (without their API keys) from the database.
func (r *appRepository) GetAllApps(ctx context.Context) ([]models.App, error) {
	_, tableSpan := otel.Tracer("db-tracer").Start(ctx, "db.query.apps")
	defer tableSpan.End()
	rows, err := r.db.Query("SELECT id, name FROM apps")
	if err != nil {
		return nil, fmt.Errorf("failed to query apps: %v", err)
	}
	defer rows.Close()

	var apps []models.App
	for rows.Next() {
		var a models.App
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, fmt.Errorf("failed to scan app row: %v", err)
		}
		apps = append(apps, a)
	}

	return apps, nil
}

// DeleteApp revokes (deletes) an app by its ID.
func (r *appRepository) DeleteApp(ctx context.Context, id int) (sql.Result, error) {
	_, tableSpan := otel.Tracer("db-tracer").Start(ctx, "db.delete.app")
	defer tableSpan.End()
	return r.db.Exec("DELETE FROM apps WHERE id = ?", id)
}

// ValidateAPIKey - Validates API Key against apps table
func (r *appRepository) ValidateAPIKey(ctx context.Context, appName, apiKey string) (string, bool, error) {
	_, span := otel.Tracer("db-tracer").Start(ctx, "db.validate_api_key")
	defer span.End()

	var storedHash string
	query := "SELECT api_key_hash FROM apps WHERE name = ?"
	err := r.db.QueryRowContext(ctx, query, appName).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			span.SetStatus(codes.Ok, "API key not found")
			return "", false, nil
		}
		span.SetStatus(codes.Error, fmt.Sprintf("db query failed: %v", err))
		span.RecordError(err)
		return "", false, err
	}

	// Compare the provided API key with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(apiKey))
	if err != nil {
		// bcrypt.CompareHashAndPassword returns an error if they don't match
		span.SetStatus(codes.Ok, "API key validation failed")
		return "", false, nil
	}

	span.SetStatus(codes.Ok, "API key validated successfully")
	return appName, true, nil
}
