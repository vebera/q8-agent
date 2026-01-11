package service

import (
	"fmt"
	"log"

	"github.com/qate/q8-agent/internal/docker"
	"github.com/qate/q8-agent/internal/domain"
	"github.com/qate/q8-agent/internal/fs"
)

// Orchestrator coordinates tenant operations
type Orchestrator struct {
	fs     *fs.Manager
	docker *docker.Runner
}

// NewOrchestrator creates a new orchestrator
func NewOrchestrator(fs *fs.Manager, docker *docker.Runner) *Orchestrator {
	return &Orchestrator{
		fs:     fs,
		docker: docker,
	}
}

// ProvisionTenant sets up a new tenant environment
func (s *Orchestrator) ProvisionTenant(req domain.TenantProvisionRequest) error {
	log.Printf("Provisioning tenant: %s (subdomain: %s)", req.ID, req.Subdomain)

	// 1. Prepare directory
	dir, err := s.fs.PrepareTenantDir(req.Subdomain)
	if err != nil {
		return fmt.Errorf("fs error: %w", err)
	}

	// 2. Write configs
	err = s.fs.WriteConfig(req.Subdomain, req.ComposeContent, req.EnvContent)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	// 3. Pull and Up
	project := fmt.Sprintf("q8-%s", req.Subdomain)

	log.Printf("Pulling images for project: %s", project)
	if out, err := s.docker.ExecuteComposePull(project, dir); err != nil {
		return fmt.Errorf("docker pull error: %s: %w", string(out), err)
	}

	log.Printf("Spinning up containers for project: %s", project)
	if out, err := s.docker.ExecuteComposeUp(project, dir); err != nil {
		return fmt.Errorf("docker up error: %s: %w", string(out), err)
	}

	log.Printf("Tenant %s provisioned successfully", req.ID)
	return nil
}

// TeardownTenant removes a tenant environment
func (s *Orchestrator) TeardownTenant(subdomain string) error {
	log.Printf("Tearing down tenant: %s", subdomain)

	project := fmt.Sprintf("q8-%s", subdomain)
	dir := s.fs.GetTenantPath(subdomain)

	// 1. Docker down
	if out, err := s.docker.ExecuteComposeDown(project, dir); err != nil {
		log.Printf("Warning: docker down failed (might already be gone): %s", string(out))
	}

	// 2. Archive files instead of removing
	newDir, err := s.fs.ArchiveTenantDir(subdomain)
	if err != nil {
		return fmt.Errorf("fs archive error: %w", err)
	}

	if newDir != "" {
		log.Printf("Tenant %s archived to %s", subdomain, newDir)
	} else {
		log.Printf("Tenant %s directory not found, nothing to archive", subdomain)
	}

	return nil
}

// RestartTenant restarts a tenant's containers
func (s *Orchestrator) RestartTenant(subdomain string) error {
	log.Printf("Restarting tenant: %s", subdomain)

	project := fmt.Sprintf("q8-%s", subdomain)
	dir := s.fs.GetTenantPath(subdomain)

	if out, err := s.docker.ExecuteComposeRestart(project, dir); err != nil {
		return fmt.Errorf("docker restart error: %s: %w", string(out), err)
	}

	return nil
}

// GetTenantStatus returns the status of a tenant's containers
func (s *Orchestrator) GetTenantStatus(subdomain string) (string, error) {
	project := fmt.Sprintf("q8-%s", subdomain)
	dir := s.fs.GetTenantPath(subdomain)

	out, err := s.docker.ExecuteComposePs(project, dir)
	if err != nil {
		return "", fmt.Errorf("docker ps error: %s: %w", string(out), err)
	}

	return string(out), nil
}

// GetTenantLogs returns the logs of a tenant's containers
func (s *Orchestrator) GetTenantLogs(subdomain string, tail int) (string, error) {
	project := fmt.Sprintf("q8-%s", subdomain)
	dir := s.fs.GetTenantPath(subdomain)

	out, err := s.docker.ExecuteComposeLogs(project, dir, tail)
	if err != nil {
		return "", fmt.Errorf("docker logs error: %s: %w", string(out), err)
	}

	return string(out), nil
}

// GetTenantImages returns the images of a tenant's containers
func (s *Orchestrator) GetTenantImages(subdomain string) (string, error) {
	project := fmt.Sprintf("q8-%s", subdomain)
	dir := s.fs.GetTenantPath(subdomain)

	out, err := s.docker.ExecuteComposeImages(project, dir)
	if err != nil {
		return "", fmt.Errorf("docker images error: %s: %w", string(out), err)
	}

	return string(out), nil
}
