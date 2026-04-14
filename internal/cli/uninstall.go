package cli

import (
	"fmt"
	"os"

	"github.com/easy-skills/easy-skills/internal/hub"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall a package and remove its components",
	Long:  "Delete all component files at their registered paths and remove Hub records",
	Run:   runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringVar(&flagName, "name", "", "Package name (required)")
	uninstallCmd.Flags().StringVar(&flagTarget, "target", "", "Target IDE: qoder or cursor (required)")
	uninstallCmd.MarkFlagRequired("name")
	uninstallCmd.MarkFlagRequired("target")
}

func runUninstall(cmd *cobra.Command, args []string) {
	pkgName := flagName
	target := flagTarget

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

	// Get components to delete their files
	components, err := h.GetComponentsByPackage(pkg.ID)
	if err != nil {
		Failf("Failed to get components: %v", err)
		return
	}

	// Delete component files
	var removed []string
	for _, comp := range components {
		if _, err := os.Stat(comp.Path); err == nil {
			if err := os.RemoveAll(comp.Path); err != nil {
				Failf("Failed to delete %s at %s: %v", comp.Name, comp.Path, err)
				return
			}
			removed = append(removed, comp.Name)
		}
	}

	// Delete installation records
	if err := h.DeleteInstallationsByPackage(pkg.ID); err != nil {
		Failf("Failed to delete installation records: %v", err)
		return
	}

	// Delete components, versions, and package from DB
	if err := h.DeletePackageFull(pkg.ID); err != nil {
		Failf("Failed to delete package records: %v", err)
		return
	}

	Success(map[string]interface{}{
		"package":  pkgName,
		"target":   target,
		"removed":  removed,
		"message":  fmt.Sprintf("Uninstalled %s: removed %d file(s) and Hub records", pkgName, len(removed)),
	})
}
