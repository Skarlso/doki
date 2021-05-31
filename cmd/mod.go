package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/Skarlso/doki/pkg/git"
	"github.com/Skarlso/doki/pkg/runner"
	"github.com/spf13/cobra"
)

var gitExtract = regexp.MustCompile(`github.com/([a-zA-Z0-9\-]+)/([a-zA-Z0-9\-]+)`)

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
)

func init() {
	RootCmd.AddCommand(modCmd)
	modCmd.AddCommand(latestCmd)
}

// runModLatestCmd .
func runModLatestCmd(cmd *cobra.Command, args []string) {
	cfg := git.Config{
		Runner: &runner.CLIRunner{},
		Token:  globalArgs.token,
	}
	provider := git.NewProvider(cfg)
	var allTags []string
	for _, imp := range args {
		if !strings.Contains(imp, "github") {
			fmt.Println("This tool does not support novel domains at the moment.")
			os.Exit(1)
		}
		if strings.Contains(imp, "@") {
			allTags = append(allTags, imp)
			continue
		}
		m := gitExtract.FindAllStringSubmatch(imp, -1)
		if len(m) == 0 {
			fmt.Printf("Failed to get latest version for import %s\n", imp)
			os.Exit(1)
		}
		if len(m[0]) < 3 {
			fmt.Println("Match does not contain repo and owner: ", m[0])
			os.Exit(1)
		}
		owner := m[0][1]
		repo := m[0][2]
		latestTag, err := provider.GetLatestRemoteTag(owner, repo)
		if err != nil {
			fmt.Printf("Failed to get latest version for import %s: %s\n", imp, err.Error())
			os.Exit(1)
		}
		allTags = append(allTags, fmt.Sprintf("%s@%s", imp, latestTag))
	}
	sort.Strings(allTags)
	fmt.Println(strings.Join(allTags, " "))
}
