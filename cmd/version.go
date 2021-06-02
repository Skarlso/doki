package cmd

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
)

var (
	// versionCmd is the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version",
		Run:   runVersionCmd,
	}
	currentVersion *semver.Version
)

// SetVersion must be called at bootstrap to pass the current build version
func SetVersion(releaseVersion string) {
	currentVersion = semver.MustParse(releaseVersion)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func runVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(currentVersion.String())
}
