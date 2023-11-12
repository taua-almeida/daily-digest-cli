package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/taua-almeida/gh-daily-digest-cli/internal"
)

func NewArgs() *internal.Args {
	return &internal.Args{
		Status:   "open",
		WithOrgs: false,
		Org:      "all",
		EnvVar:   "GITHUB_TOKEN",
		RepoName: "all",
	}
}

var rootCmd = &cobra.Command{
	Use:   "gh-digest",
	Short: "A daily digest of your github work",
	Long:  "A daily digest of your github tasks: pending reviews, PR status, notifications, commits, comments",
	Args:  cobra.NoArgs,
	RunE:  runRootCmd,
}

func parseCommandArgs(cmd *cobra.Command) (*internal.Args, error) {
	args := NewArgs()

	var err error
	args.EnvVar, err = cmd.Flags().GetString("env-var")
	if err != nil {
		return nil, err
	}

	args.RepoName, err = cmd.Flags().GetString("repo")
	if err != nil {
		return nil, err
	}

	args.Status, err = cmd.Flags().GetString("status")
	if err != nil {
		return nil, err
	}

	args.WithOrgs, err = cmd.Flags().GetBool("with-orgs")
	if err != nil {
		return nil, err
	}

	args.Org, err = cmd.Flags().GetString("org")
	if err != nil {
		return nil, err
	}

	return args, nil
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	receivedCommandArgs, err := parseCommandArgs(cmd)
	if err != nil {
		return err
	}

	if err := validateCommandArgs(receivedCommandArgs); err != nil {
		return err
	}

	pullRequests, err := internal.FetchPullRequests(receivedCommandArgs)
	if err != nil {
		return err
	}

	internal.PrintPullRequests(pullRequests)

	return nil

}

func Execute() error {
	return rootCmd.Execute()
}

func validateCommandArgs(args *internal.Args) error {
	err := isValidStatusArg(args.Status)
	if err != nil {
		return err
	}

	if !args.WithOrgs && args.Org != "" {
		return errors.New("you can't use the --org flag without --with-orgs")
	}

	return nil
}

func isValidStatusArg(val string) error {
	allowedStatuses := []string{"open", "closed", "all"}
	for _, status := range allowedStatuses {
		if val == status {
			return nil
		}
	}
	return fmt.Errorf("unsupported status: %s, allowed statuses are: %s", val, strings.Join(allowedStatuses, ", "))
}

func init() {
	commandArgs := NewArgs()
	rootCmd.Flags().StringVarP(&commandArgs.EnvVar, "env-var", "", commandArgs.EnvVar, "The name of the environment variable to fetch github's token")
	rootCmd.Flags().StringVarP(&commandArgs.RepoName, "repo", "r", commandArgs.RepoName, "The name of the repository to fetch pull requests from [all|<repo_name>]")
	rootCmd.Flags().StringVarP(&commandArgs.Status, "status", "s", commandArgs.Status, "The status of the pull requests to be fetched [open|closed|all]")
	rootCmd.Flags().BoolVarP(&commandArgs.WithOrgs, "with-orgs", "w", commandArgs.WithOrgs, "Fetch pull requests from organizations the user is part of")
	rootCmd.Flags().StringVarP(&commandArgs.Org, "org", "o", commandArgs.Org, "The name of the organization to fetch pull requests from")
}
