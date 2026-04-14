package hub

import (
	"github.com/easy-skills/easy-skills/pkg/types"
)

// GetInstallationsByPackage returns all installations for a package
func (h *Hub) GetInstallationsByPackage(packageID string) ([]*types.Installation, error) {
	rows, err := h.db.Query(`
		SELECT i.id, i.component_id, i.ide, i.scope, i.install_path, i.installed_at
		FROM installations i
		JOIN components c ON i.component_id = c.id
		WHERE c.package_id = ?
	`, packageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var installations []*types.Installation
	for rows.Next() {
		inst := &types.Installation{}
		if err := rows.Scan(&inst.ID, &inst.ComponentID, &inst.IDE, &inst.Scope, &inst.InstallPath, &inst.InstalledAt); err != nil {
			return nil, err
		}
		installations = append(installations, inst)
	}
	return installations, nil
}

// DeleteInstallationsByPackage deletes all installations for a package
func (h *Hub) DeleteInstallationsByPackage(packageID string) error {
	_, err := h.db.Exec(`
		DELETE FROM installations 
		WHERE component_id IN (
			SELECT id FROM components WHERE package_id = ?
		)
	`, packageID)
	return err
}

// HasInstallations checks if a package has any installations
func (h *Hub) HasInstallations(packageID string) (bool, error) {
	var count int
	err := h.db.QueryRow(`
		SELECT COUNT(*) FROM installations 
		WHERE component_id IN (
			SELECT id FROM components WHERE package_id = ?
		)
	`, packageID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetComponentsByPackage returns all components for a package
func (h *Hub) GetComponentsByPackage(packageID string) ([]*types.Component, error) {
	rows, err := h.db.Query(`
		SELECT c.id, c.package_id, c.package_version_id, c.name, c.type, c.path
		FROM components c
		WHERE c.package_id = ?
	`, packageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []*types.Component
	for rows.Next() {
		c := &types.Component{}
		if err := rows.Scan(&c.ID, &c.PackageID, &c.PackageVersionID, &c.Name, &c.Type, &c.Path); err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return components, nil
}
