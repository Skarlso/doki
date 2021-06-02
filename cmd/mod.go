package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Skarlso/doki/pkg/git"
	"github.com/Skarlso/doki/pkg/gomod"
	"github.com/Skarlso/doki/pkg/runner"
)

var (
	// ModCmd is the command root for 'mod ...' commands.
	modCmd = &cobra.Command{
		Use:   "mod",
		Short: "Manipulate the go mod file.",
	}
	latestCmd = &cobra.Command{
		Use:   "latest",
		Short: "Get the latest version for a given module.",
		Run:   runModLatestCmd,
	}
	replaceCmd = &cobra.Command{
		Use:   "replace",
		Short: "Replace a module with a designated target",
		Run:   runModReplaceCmd,
	}
	replaceArgs struct {
		replacements []string
	}
)

func init() {
	RootCmd.AddCommand(modCmd)
	modCmd.AddCommand(latestCmd)
	modCmd.AddCommand(replaceCmd)

	f := replaceCmd.PersistentFlags()
	f.StringSliceVarP(&replaceArgs.replacements, "replacements", "p", nil, "List of replacements (--replace mod=mod@ver, mod2=mod2@ver).")
}

func runModReplaceCmd(cmd *cobra.Command, args []string) {
	result := make([]string, 0)
	for _, replacement := range replaceArgs.replacements {
		result = append(result, fmt.Sprintf("-replace %s", replacement))
	}
	fmt.Print(strings.Join(result, " "))
}

// runModLatestCmd .
func runModLatestCmd(cmd *cobra.Command, args []string) {
	cfg := git.Config{
		Runner: &runner.CLIRunner{},
		Token:  globalArgs.token,
	}
	provider := git.NewProvider(cfg)
	modProvider := gomod.NewProvider(gomod.Config{
		VCS: provider,
	})
	allTags, err := modProvider.GetLatestModuleList(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(strings.Join(allTags, " "))
}
