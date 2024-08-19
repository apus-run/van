package upgrade

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the van tools",
	Long:  "Upgrade the van tools. Example: van upgrade",
	Run:   Run,
}

// Run upgrade the van tools.
func Run(cmd *cobra.Command, args []string) {

}
