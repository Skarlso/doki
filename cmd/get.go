package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Skarlso/doki/pkg/git"
	"github.com/Skarlso/doki/pkg/runner"
)

var (
	// GetCmd is the command root for 'get ...' commands.
	getDevTagCmd = &cobra.Command{
		Use:   "get-dev-tag",
		Short: "Get information about the current repository.",
		Run:   runTagCmd,
	}
)

func init() {
	RootCmd.AddCommand(getDevTagCmd)
}

// runTagCmd run dev tag command.
func runTagCmd(cmd *cobra.Command, args []string) {
	cfg := git.Config{
		Runner: &runner.CLIRunner{},
		Token:  globalArgs.token,
	}
	provider := git.NewProvider(cfg)
	tag, err := provider.GetDevTag()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(tag)
}
