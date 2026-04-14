package cli

import "github.com/spf13/cobra"

var (
	flagName       string
	flagTarget     string
	flagComponents []string
	flagIDE        string
	flagScope      string
	flagPath       string
)

func bindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&flagName, "name", "", "Package name")
	cmd.Flags().StringVar(&flagTarget, "target", "", "Target IDE (qoder/cursor)")
	cmd.Flags().StringVar(&flagIDE, "ide", "", "IDE to install to")
	cmd.Flags().StringVar(&flagScope, "scope", "user", "Scope (user/project)")
}
