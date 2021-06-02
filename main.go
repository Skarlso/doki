package main

import (
	"fmt"
	"os"

	"github.com/Skarlso/doki/cmd"
)

var (
	version = "v0.0.0-dev"
)

func init() {
	cmd.SetVersion(version)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println("Error while running root command: ", err)
		os.Exit(1)
	}
}
