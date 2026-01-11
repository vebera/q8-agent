package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/qate/q8-agent/internal/domain"
	"github.com/qate/q8-agent/internal/service"
)

// Handler handles API requests
type Handler struct {
	service *service.Orchestrator
}

// NewHandler creates a new API handler
func NewHandler(s *service.Orchestrator) *Handler {
	return &Handler{service: s}
}

// Provision handles tenant provisioning
func (h *Handler) Provision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.TenantProvisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" || req.Subdomain == "" {
		http.Error(w, "Missing required fields (id, subdomain)", http.StatusBadRequest)
		return
	}

	if err := h.service.ProvisionTenant(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "provisioned", "id": req.ID})
}

// Teardown handles tenant teardown
func (h *Handler) Teardown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simple path parsing /v1/tenants/teardown/{subdomain}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Missing subdomain", http.StatusBadRequest)
		return
	}
	subdomain := parts[len(parts)-1]

	if err := h.service.TeardownTenant(subdomain); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "torn_down", "subdomain": subdomain})
}

// Restart handles tenant restart
func (h *Handler) Restart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Missing subdomain", http.StatusBadRequest)
		return
	}
	subdomain := parts[len(parts)-1]

	if err := h.service.RestartTenant(subdomain); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "restarted", "subdomain": subdomain})
}

// Status handles tenant status request
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Missing subdomain", http.StatusBadRequest)
		return
	}
	subdomain := parts[len(parts)-1]

	status, err := h.service.GetTenantStatus(subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status))
}

// Logs handles tenant logs request
func (h *Handler) Logs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Missing subdomain", http.StatusBadRequest)
		return
	}
	subdomain := parts[len(parts)-1]

	// Simple tail parsing from query param
	tail := 100
	if t := r.URL.Query().Get("tail"); t != "" {
		fmt.Sscanf(t, "%d", &tail)
	}

	logs, err := h.service.GetTenantLogs(subdomain, tail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(logs))
}

// Images handles tenant images request
func (h *Handler) Images(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Missing subdomain", http.StatusBadRequest)
		return
	}
	subdomain := parts[len(parts)-1]

	images, err := h.service.GetTenantImages(subdomain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(images))
}

// CreateDatabase handles mongo database/user creation
func (h *Handler) CreateDatabase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req domain.MongoDBUserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Host == "" || req.AdminUser == "" || req.DatabaseName == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateMongoDBUser(req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "database_configured"})
}
