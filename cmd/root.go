package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// RootCmd is the root command of this tool
	RootCmd = &cobra.Command{
		Use:   "doki",
		Short: "Weaveworks build and version syncing tool.",
		Run:   RunUsage,
	}
)

// RunUsage show usage
func RunUsage(cmd *cobra.Command, args []string) {
	if err := cmd.Usage(); err != nil {
		fmt.Println("Error showing usage: ", err.Error())
		os.Exit(1)
	}
}
