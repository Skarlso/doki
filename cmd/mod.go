package cmd

import (
	"fmt"
	"os"

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
	updateArgs struct {
		listRepositoryShorts bool
	}
)

func init() {
	RootCmd.AddCommand(modCmd)
	modCmd.AddCommand(updateCmd)

	f := updateCmd.PersistentFlags()
	// Persistent flags
	f.BoolVarP(&updateArgs.listRepositoryShorts, "repositories", "r", false, "List all available repository shorts for update command.")
}

// runModUpdateCmd .
func runModUpdateCmd(cmd *cobra.Command, args []string) {
	if updateArgs.listRepositoryShorts {
		fmt.Println("The following repository shorts are available:")
		for k, v := range repositories {
			fmt.Printf("	%s = %s\n", k, v)
		}
		os.Exit(0)
	}
}
