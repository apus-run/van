package create

import (
	"github.com/spf13/cobra"

	"github.com/apus-run/van/cmd/van/internal/create/model"
	"github.com/apus-run/van/cmd/van/internal/create/service"
)

// Cmd represents the new command.
var Cmd = &cobra.Command{
	Use:   "new",
	Short: "create new",
	Long:  "Generate the new files.",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
	Args: cobra.NoArgs,
}

func init() {
	Cmd.AddCommand(model.Cmd)
	Cmd.AddCommand(service.Cmd)
}
