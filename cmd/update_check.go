package cmd

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"

	"github.com/Skarlso/doki/pkg/git"
)

var (
	// updateCheckCmd is the root for all go related commands
	updateCheckCmd = &cobra.Command{
		Use:   "update-check",
		Short: "Check if this version is the latest version",
		Run:   runUpdateCheckCmd,
	}
)

func init() {
	RootCmd.AddCommand(updateCheckCmd)
}

// Run the service
func runUpdateCheckCmd(cmd *cobra.Command, args []string) {
	provider := git.NewProvider()
	latestVersion, err := provider.GetLatestRemoteTag()
	if err != nil {
		fmt.Println("Failed to get latest version: ", err)
		os.Exit(1)
	}
	lv := semver.MustParse(latestVersion)

	if !lv.Equal(currentVersion) {
		fmt.Printf("Doki is not the latest version (%s -> %s)\n", lv.String(), latestVersion)
		fmt.Println()
		fmt.Println("To update, run:")
		fmt.Println()
		fmt.Println("  GOSUMDB=off GOPROXY=direct GO111MODULE=on go get github.com/Skarlso/doki@v" + latestVersion)
		fmt.Println()
		os.Exit(1)
	}
	fmt.Println("Dōki is up to date")
}
