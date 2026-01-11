package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/qate/q8-agent/internal/api"
	"github.com/qate/q8-agent/internal/config"
	"github.com/qate/q8-agent/internal/docker"
	"github.com/qate/q8-agent/internal/fs"
	"github.com/qate/q8-agent/internal/service"
)

func main() {
	// 1. Load config
	cfg := config.LoadConfig()

	// 2. Initialize components
	fsManager := fs.NewManager(cfg.TenantsRoot)
	dockerRunner := docker.NewRunner()

	// Check if docker is available
	if !dockerRunner.IsInstalled() {
		log.Fatal("Fatal: docker compose is not installed or accessible")
	}

	orchestrator := service.NewOrchestrator(cfg, fsManager, dockerRunner)
	handler := api.NewHandler(orchestrator)

	// 3. Setup Routes
	mux := http.NewServeMux()

	// Add routes with Auth Middleware
	mux.HandleFunc("/v1/tenants/provision", api.AuthMiddleware(cfg, handler.Provision))
	mux.HandleFunc("/v1/tenants/teardown/", api.AuthMiddleware(cfg, handler.Teardown))
	mux.HandleFunc("/v1/tenants/restart/", api.AuthMiddleware(cfg, handler.Restart))
	mux.HandleFunc("/v1/tenants/status/", api.AuthMiddleware(cfg, handler.Status))
	mux.HandleFunc("/v1/tenants/logs/", api.AuthMiddleware(cfg, handler.Logs))
	mux.HandleFunc("/v1/tenants/images/", api.AuthMiddleware(cfg, handler.Images))

	// Health check (no auth)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	// 4. Start Server
	log.Printf("Q8 Agent starting on port %s...", cfg.Port)
	log.Printf("Tenants root: %s", cfg.TenantsRoot)
	printRoutes()

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

func printRoutes() {
	log.Println("Supported API Methods:")
	log.Println("  [POST] /v1/tenants/provision  - Provision a new tenant environment")
	log.Println("  [POST] /v1/tenants/teardown/  - Remove a tenant environment")
	log.Println("  [POST] /v1/tenants/restart/   - Restart tenant containers")
	log.Println("  [GET]  /v1/tenants/status/    - Get container status")
	log.Println("  [GET]  /v1/tenants/logs/      - Get container logs")
	log.Println("  [GET]  /v1/tenants/images/    - Get container image information")
	log.Println("  [GET]  /health                - Agent health check")
}
