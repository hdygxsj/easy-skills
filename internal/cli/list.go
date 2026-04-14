package cli

import (
	"encoding/json"
	"os"

	"github.com/easy-skills/easy-skills/internal/hub"
	"github.com/spf13/cobra"
)

var listJSON bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages in Hub",
	Long:  "List all packages in the local Hub that have installation records, optionally filtered by target",
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&flagTarget, "target", "", "Filter by target (qoder/cursor)")
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
}

func runList(cmd *cobra.Command, args []string) {
	h, err := hub.NewHub()
	if err != nil {
		Failf("Failed to create hub: %v", err)
		return
	}
	defer h.Close()

	packages, err := h.ListPackages(flagTarget)
	if err != nil {
		Failf("Failed to list packages: %v", err)
		return
	}

	// Enrich with version info and installation status
	var result []map[string]interface{}
	for _, pkg := range packages {
		// Only show packages that have installation records
		hasInstallations, err := h.HasInstallations(pkg.ID)
		if err != nil {
			continue
		}
		if !hasInstallations {
			continue
		}

		version, _ := h.GetCurrentVersion(pkg.ID)
		components, _ := h.GetComponentsByVersion(version.ID)
		installations, _ := h.GetInstallationsByPackage(pkg.ID)

		result = append(result, map[string]interface{}{
			"package":          pkg,
			"current_version":  version,
			"components":       components,
			"components_count": len(components),
			"installations":    installations,
			"installed":        true,
		})
	}

	if listJSON {
		// Raw JSON output for API
		json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
			"packages": result,
			"count":    len(result),
		})
	} else {
		Success(map[string]interface{}{
			"packages": result,
			"count":    len(result),
		})
	}
}
