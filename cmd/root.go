package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	envKeyPrefix = "DOKI_"
)

var (
	// RootCmd is the root command of this tool
	RootCmd = &cobra.Command{
		Use:   "doki",
		Short: "Weaveworks build and version syncing tool.",
		Run:   RunUsage,
	}
	globalArgs struct {
		token string
	}
)

func init() {
	token := getEnvOrDefault("TOKEN", "")
	RootCmd.PersistentFlags().StringVarP(&globalArgs.token, "token", "t", token, "Optional GitHub token to get releases for non-public or restricted repositories.")
}

// RunUsage show usage
func RunUsage(cmd *cobra.Command, args []string) {
	if err := cmd.Usage(); err != nil {
		fmt.Println("Error showing usage: ", err.Error())
		os.Exit(1)
	}
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(envKeyPrefix + key); v != "" {
		return v
	}
	return def
}
