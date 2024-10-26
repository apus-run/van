package model

import (
	"github.com/spf13/cobra"
)

// Cmd represents the new command.
var Cmd = &cobra.Command{
	Use:   "new",
	Short: "Create a model",
	Long:  "Create a model using the repository template. Example: van new helloworld",
	Run:   run,
	Args:  cobra.NoArgs,
}

func run(cmd *cobra.Command, args []string) {

}
