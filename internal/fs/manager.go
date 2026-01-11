package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager handles file system operations for tenants
type Manager struct {
	root string
}

// NewManager creates a new file system manager
func NewManager(root string) *Manager {
	return &Manager{root: root}
}

// PrepareTenantDir creates the tenant directory and returns its path
func (m *Manager) PrepareTenantDir(subdomain string) (string, error) {
	path := filepath.Join(m.root, subdomain)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create tenant directory: %w", err)
	}
	return path, nil
}

// WriteConfig write the docker-compose and .env files
func (m *Manager) WriteConfig(subdomain, compose, env string) error {
	dir := filepath.Join(m.root, subdomain)

	err := os.WriteFile(filepath.Join(dir, "docker-compose.yml"), []byte(compose), 0644)
	if err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}

	err = os.WriteFile(filepath.Join(dir, ".env"), []byte(env), 0644)
	if err != nil {
		return fmt.Errorf("failed to write .env: %w", err)
	}

	return nil
}

// RemoveTenantDir deletes the tenant directory
func (m *Manager) RemoveTenantDir(subdomain string) error {
	path := m.GetTenantPath(subdomain)
	return os.RemoveAll(path)
}

// GetTenantPath returns the absolute path for a tenant
func (m *Manager) GetTenantPath(subdomain string) string {
	return filepath.Join(m.root, subdomain)
}
