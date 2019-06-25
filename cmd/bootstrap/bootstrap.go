package bootstrap

import (
	"fmt"

	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/environment"
)

// PluginSettingsID is used as a key in the root settings map of plugin settings
const PluginSettingsID = "bootstrap"

// NewBootstrapCommand creates a new command for bootstrapping airshipctl
func NewBootstrapCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	bootstrapRootCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstraps airshipctl",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Under construction\n")
		},
	}

	return bootstrapRootCmd
}
