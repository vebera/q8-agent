package fs

import (
	"crypto/rand"
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

// ArchiveTenantDir renames the tenant directory with a UUID suffix
func (m *Manager) ArchiveTenantDir(subdomain string) (string, error) {
	oldPath := m.GetTenantPath(subdomain)

	// Check if directory exists
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return "", nil // Nothing to archive
	}

	// Generate UUID
	uuid, err := newUUID()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	newDirName := fmt.Sprintf("%s-%s", subdomain, uuid)
	newPath := filepath.Join(m.root, newDirName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return "", fmt.Errorf("failed to archive directory: %w", err)
	}

	return newDirName, nil
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

// newUUID generates a random UUID (version 4)
func newUUID() (string, error) {
	var u [16]byte
	_, err := rand.Read(u[:])
	if err != nil {
		return "", err
	}
	u[6] = (u[6] & 0x0f) | 0x40 // Version 4
	u[8] = (u[8] & 0x3f) | 0x80 // Variant is 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}
