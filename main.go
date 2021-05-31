package main

import (
	"fmt"
	"os"

	"github.com/Skarlso/doki/cmd"
)

func init() {
	cmd.SetVersion(releaseVersion)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println("Error while running root command: ", err)
		os.Exit(1)
	}
}
