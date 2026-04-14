package cli

import (
	"github.com/spf13/cobra"
)

var reinstallCmd = &cobra.Command{
	Use:   "reinstall",
	Short: "Verify a package's components (alias for install)",
	Long:  "Check that all registered components exist at their paths. Same as install.",
	Run:   runReinstall,
}

func init() {
	rootCmd.AddCommand(reinstallCmd)
	reinstallCmd.Flags().StringVar(&flagName, "name", "", "Package name (required)")
	reinstallCmd.Flags().StringVar(&flagTarget, "target", "", "Target IDE: qoder or cursor (required)")
	reinstallCmd.Flags().StringVar(&flagIDE, "ide", "", "IDE (qoder/cursor)")
	reinstallCmd.Flags().StringVar(&flagScope, "scope", "user", "Scope: user or project")
	reinstallCmd.MarkFlagRequired("name")
	reinstallCmd.MarkFlagRequired("target")
}

func runReinstall(cmd *cobra.Command, args []string) {
	// In the new model, reinstall is the same as install (verify)
	runInstall(cmd, args)
}
