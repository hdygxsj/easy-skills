package cli

import (
	"fmt"

	"github.com/easy-skills/easy-skills/internal/adapters/impl"
	"github.com/easy-skills/easy-skills/internal/hub"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall a package from IDE",
	Long:  "Uninstall a package's components from the target IDE and remove installation records",
	Run:   runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringVar(&flagName, "name", "", "Package name (required)")
	uninstallCmd.Flags().StringVar(&flagTarget, "target", "", "Target IDE: qoder or cursor (required)")
	uninstallCmd.Flags().StringVar(&flagIDE, "ide", "", "IDE to uninstall from (qoder/cursor)")
	uninstallCmd.Flags().StringVar(&installScope, "scope", "user", "Installation scope: user or project")
	uninstallCmd.MarkFlagRequired("name")
	uninstallCmd.MarkFlagRequired("target")
}

func runUninstall(cmd *cobra.Command, args []string) {
	pkgName := flagName
	target := flagTarget
	ide := flagIDE
	if ide == "" {
		ide = target
	}
	scope := installScope

	// Get IDE adapter
	adapter, err := impl.Get(ide)
	if err != nil {
		Failf("Unsupported IDE: %s", ide)
		return
	}

	// Get installed path
	path, err := adapter.GetInstalledPath(pkgName, scope)
	if err != nil {
		Failf("Failed to get path: %v", err)
		return
	}

	// Uninstall from IDE (delete files)
	if err := adapter.UninstallComponent(pkgName, path); err != nil {
		Failf("Failed to uninstall from IDE: %v", err)
		return
	}

	// Remove installation records from Hub database
	h, err := hub.NewHub()
	if err != nil {
		Failf("Failed to create hub: %v", err)
		return
	}
	defer h.Close()

	// Get package
	pkg, err := h.GetPackage(pkgName, target)
	if err != nil {
		// If package doesn't exist, that's OK - just return success
		Success(map[string]interface{}{
			"package": pkgName,
			"scope":   scope,
			"ide":     ide,
			"message": fmt.Sprintf("Uninstalled %s from IDE (no Hub records found)", pkgName),
		})
		return
	}

	// Delete installation records for this package
	if err := h.DeleteInstallationsByPackage(pkg.ID); err != nil {
		Failf("Failed to delete installation records: %v", err)
		return
	}

	Success(map[string]interface{}{
		"package": pkgName,
		"scope":   scope,
		"ide":     ide,
		"message": fmt.Sprintf("Uninstalled %s and removed Hub records", pkgName),
	})
}
