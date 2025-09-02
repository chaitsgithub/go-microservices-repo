package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"chaits.org/go-microservices-repo/internal/models"
	"chaits.org/go-microservices-repo/internal/repositories"
	"chaits.org/go-microservices-repo/pkg/general/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type AppsHandler struct {
	appRepo repositories.AppRepository
}

func NewAppsHandler(db *repositories.DBManager) *AppsHandler {
	return &AppsHandler{
		appRepo: db.AppRepo,
	}
}

func (a *AppsHandler) GetAppsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("handler.name", "getapps-handler"))

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	dbCtx, dbSpan := otel.Tracer("db-tracer").Start(ctx, "db.query")
	defer dbSpan.End()

	apps, err := a.appRepo.GetAllApps(dbCtx)
	if err != nil {
		logger.Logger.WithError(err).Error("Error getting apps from db")
		dbSpan.SetStatus(codes.Error, "db query failed")
		http.Error(w, "Error getting apps from db", http.StatusInternalServerError)
		return
	}

	dbSpan.SetStatus(codes.Ok, "db query successful")
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(apps); err != nil {
		http.Error(w, "Error formatting response", http.StatusInternalServerError)
		return
	}
}

// generateAPIKey creates a new, cryptographically secure random API key.
// It returns the key as a hexadecimal string.
func generateAPIKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func (a *AppsHandler) RegisterAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("handler.name", "registerapp-handler"))

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	dbCtx, dbSpan := otel.Tracer("db-tracer").Start(ctx, "db.insert")
	defer dbSpan.End()

	var req createAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error with request body", http.StatusBadRequest)
		span.SetStatus(codes.Error, "Error with request body")
		return
	}

	// Validate required fields. We now only expect the app name.
	if req.Name == "" {
		http.Error(w, "Bad request: name is required", http.StatusBadRequest)
		span.SetStatus(codes.Error, "Error with request body")
		return
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		http.Error(w, "Error generating API Key", http.StatusInternalServerError)
		span.SetStatus(codes.Error, "Error generating API Key")
		return
	}

	newApp := &models.App{Name: req.Name, APIKey: apiKey}

	app, err := a.appRepo.CreateApp(dbCtx, newApp)
	if err != nil {
		logger.Logger.WithError(err).Error("Error inserting app into db")
		dbSpan.SetStatus(codes.Error, "db query failed")
		http.Error(w, "failed to create app", http.StatusInternalServerError)
		return
	}

	dbSpan.SetStatus(codes.Ok, "db insert successful")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)

}

func (a *AppsHandler) RevokeAppHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("handler.name", "revokeapp-handler"))
	fmt.Println("In Revoke Handler")

	// Only allow DELETE requests for this handler.
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the ID from the URL query parameters.
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	// Convert the ID string to an integer.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	dbCtx, dbSpan := otel.Tracer("db-tracer").Start(ctx, "db.delete")
	defer dbSpan.End()

	res, err := a.appRepo.DeleteApp(dbCtx, id)
	if err != nil {
		logger.Logger.WithError(err).Errorf("Error deleting app from db for id : %d", id)
		dbSpan.SetStatus(codes.Error, "db delete failed")
		http.Error(w, fmt.Sprintf("error deleting app from db for id : %d", id), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		dbSpan.SetStatus(codes.Error, "db delete failed")
		http.Error(w, "App not found or could not be deleted", http.StatusNotFound)
		dbSpan.SetStatus(codes.Ok, "app not found or could not be deleted")
		return
	}

	dbSpan.SetStatus(codes.Ok, "db delete successful")
	dbSpan.End()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("App deleted successfully! id: %d", id)))
}
