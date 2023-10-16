package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "daily-digest",
	Short: "A daily digest of your pending tasks",
	Long:  "A daily digest of your pending tasks: github, jira, etc",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(gitHubCmd)
}
