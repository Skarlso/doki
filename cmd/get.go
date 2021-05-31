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
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get information about the current repository.",
	}
	// GetCmd is the command root for 'get ...' commands.
	devCmd = &cobra.Command{
		Use:   "dev",
		Short: "Get information about the current repository.",
	}
	// tagCmd is the command root for 'get ...' commands.
	tagCmd = &cobra.Command{
		Use:   "tag",
		Short: "Get information about the current repository.",
		Run:   runTagCmd,
	}
)

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.AddCommand(devCmd)
	devCmd.AddCommand(tagCmd)
}

// runTagCmd .
func runTagCmd(cmd *cobra.Command, args []string) {
	cfg := git.Config{
		Runner: &runner.CLIRunner{},
		Token:  globalArgs.token,
	}
	provider := git.NewProvider(cfg)
	owner, repo, err := provider.GetOwnerAndRepoFromLocal()
	if err != nil {
		fmt.Println("Failed to get info from local repository: ", err)
		os.Exit(1)
	}
	latestRelease, err := provider.GetLatestRemoteTag(owner, repo)
	if err != nil {
		fmt.Println("Failed to get latest version for dev tagging.")
		os.Exit(1)
	}
	branch, err := provider.GetCurrentBranch()
	if err != nil {
		fmt.Println("Failed to get current branch.")
		os.Exit(1)
	}
	fmt.Printf("%s-%s\n", latestRelease, branch)
}
