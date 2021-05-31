package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// ModCmd is the command root for 'mod ...' commands.
	modCmd = &cobra.Command{
		Use:   "mod",
		Short: "Manipulate the go mod file.",
	}
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update the modules file",
		Run:   runModUpdateCmd,
	}
)

func init() {
	RootCmd.AddCommand(modCmd)
	modCmd.AddCommand(updateCmd)
}

// runModUpdateCmd .
func runModUpdateCmd(cmd *cobra.Command, args []string) {

}
