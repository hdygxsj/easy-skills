package cli

import (
	"time"

	"github.com/easy-skills/easy-skills/internal/hub"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a package to Hub",
	Long:  "Register a package from source URL to the local Hub",
	Run:   runRegister,
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVar(&flagName, "name", "", "Package name (required)")
	registerCmd.Flags().StringVar(&flagTarget, "target", "", "Target IDE: qoder or cursor (required)")
	registerCmd.Flags().StringVar(&flagSource, "source", "", "Source URL (optional)")
	registerCmd.MarkFlagRequired("name")
	registerCmd.MarkFlagRequired("target")
}

func runRegister(cmd *cobra.Command, args []string) {
	h, err := hub.NewHub()
	if err != nil {
		Failf("Failed to create hub: %v", err)
		return
	}
	defer h.Close()

	// Check if package already exists
	exists, err := h.PackageExists(flagName, flagTarget)
	if err != nil {
		Failf("Failed to check package: %v", err)
		return
	}
	if exists {
		Fail("Package already exists. Use 'upgrade' to update it.")
		return
	}

	// Create package
	pkg, err := h.CreatePackage(flagName, flagTarget, flagSource, "{}")
	if err != nil {
		Failf("Failed to create package: %v", err)
		return
	}

	// Create initial version
	version, err := h.CreateVersion(pkg.ID, timestampVersion(), "{}")
	if err != nil {
		Failf("Failed to create version: %v", err)
		return
	}

	Success(map[string]interface{}{
		"package": pkg,
		"version": version,
		"message": "Package registered successfully",
	})
}

func timestampVersion() int64 {
	return time.Now().Unix()
}
