package cli

import (
	"github.com/easy-skills/easy-skills/internal/adapters/impl"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check installation status",
	Long:  "Check which packages are installed in each IDE",
	Run:   runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVar(&flagIDE, "ide", "", "Filter by IDE (qoder/cursor)")
}

func runStatus(cmd *cobra.Command, args []string) {
	result := map[string]interface{}{
		"user":    map[string]interface{}{},
		"project": map[string]interface{}{},
	}

	ides := []string{"qoder", "cursor"}
	if flagIDE != "" {
		ides = []string{flagIDE}
	}

	for _, ideName := range ides {
		adapter, err := impl.Get(ideName)
		if err != nil {
			continue
		}

		// Check user scope
		userComponents, _ := adapter.ListInstalled("user")
		result["user"].(map[string]interface{})[ideName] = map[string]interface{}{
			"components": userComponents,
			"count":      len(userComponents),
		}

		// Check project scope
		projectComponents, _ := adapter.ListInstalled("project")
		result["project"].(map[string]interface{})[ideName] = map[string]interface{}{
			"components": projectComponents,
			"count":      len(projectComponents),
		}
	}

	Success(result)
}
