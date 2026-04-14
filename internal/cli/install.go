package cli

import (
	"fmt"
	"os"

	"github.com/easy-skills/easy-skills/internal/hub"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Verify installed components of a package",
	Long: `Check that all registered components of a package exist at their paths.

Since components are registered after installation, this command verifies
that all component files are still in place and reports their status.`,
	Run: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVar(&flagName, "name", "", "Package name (required)")
	installCmd.Flags().StringVar(&flagTarget, "target", "", "Target IDE: qoder or cursor (required)")
	installCmd.Flags().StringVar(&flagIDE, "ide", "", "IDE (qoder/cursor)")
	installCmd.Flags().StringVar(&flagScope, "scope", "user", "Scope: user or project")
	installCmd.MarkFlagRequired("name")
	installCmd.MarkFlagRequired("target")
}

func runInstall(cmd *cobra.Command, args []string) {
	pkgName := flagName
	target := flagTarget

	// Get package from Hub
	h, err := hub.NewHub()
	if err != nil {
		Failf("Failed to create hub: %v", err)
		return
	}
	defer h.Close()

	pkg, err := h.GetPackage(pkgName, target)
	if err != nil {
		Failf("Package not found: %s/%s", pkgName, target)
		return
	}

	// Get current version
	version, err := h.GetCurrentVersion(pkg.ID)
	if err != nil {
		Failf("Failed to get version: %v", err)
		return
	}

	// Get components
	components, err := h.GetComponentsByVersion(version.ID)
	if err != nil {
		Failf("Failed to get components: %v", err)
		return
	}

	if len(components) == 0 {
		Fail("Package has no components")
		return
	}

	// Check each component's path exists
	var results []map[string]interface{}
	allOK := true
	for _, comp := range components {
		exists := false
		if _, err := os.Stat(comp.Path); err == nil {
			exists = true
		}
		if !exists {
			allOK = false
		}
		results = append(results, map[string]interface{}{
			"name":   comp.Name,
			"type":   comp.Type,
			"path":   comp.Path,
			"exists": exists,
		})
	}

	status := "ok"
	msg := fmt.Sprintf("All %d component(s) verified", len(components))
	if !allOK {
		status = "missing"
		msg = "Some components are missing from their registered paths"
	}

	Success(map[string]interface{}{
		"package":    pkg.Name,
		"target":     pkg.Target,
		"status":     status,
		"components": results,
		"message":    msg,
	})
}
