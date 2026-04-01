package cli

import (
	"github.com/easy-skills/easy-skills/internal/hub"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show package details",
	Long:  "Show detailed information about a package including its components",
	Run:   runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVar(&flagName, "name", "", "Package name (required)")
	infoCmd.Flags().StringVar(&flagTarget, "target", "", "Target IDE (required)")
	infoCmd.MarkFlagRequired("name")
	infoCmd.MarkFlagRequired("target")
}

func runInfo(cmd *cobra.Command, args []string) {
	h, err := hub.NewHub()
	if err != nil {
		Failf("Failed to create hub: %v", err)
		return
	}
	defer h.Close()

	pkg, err := h.GetPackage(flagName, flagTarget)
	if err != nil {
		Failf("Package not found: %s/%s", flagName, flagTarget)
		return
	}

	version, _ := h.GetCurrentVersion(pkg.ID)
	components, _ := h.GetComponentsByVersion(version.ID)

	Success(map[string]interface{}{
		"package":    pkg,
		"version":    version,
		"components": components,
	})
}
