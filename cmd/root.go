package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gitpilot",
	Short: "GitPilot: A CLI tool to enhance your Git workflow",
	Long:  `GitPilot is a CLI tool designed to help you manage Git branches, track remote status, and handle merge conflicts more effectively.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to GitPilot!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
