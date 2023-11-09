package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/taua-almeida/daily-digest-cli/internal"
)

var (
	// Used for flags validation.
	status   string
	withOrgs bool
	org      string
)

var rootCmd = &cobra.Command{
	Use:   "gh-digest",
	Short: "A daily digest of your github work",
	Long:  "A daily digest of your github tasks: pending reviews, PR status, notifications, commits, comments",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateCommandArgs()

		if err != nil {
			return err
		}

		envVar, err := cmd.Flags().GetString("env-var")
		if err != nil {
			return err
		}
		gitToken := internal.GetGitToken(envVar)

		repo, err := cmd.Flags().GetString("repo")
		if err != nil {
			return err
		}

		status, err := cmd.Flags().GetString("status")
		if err != nil {
			return err
		}

		withOrgs, err := cmd.Flags().GetBool("with-orgs")
		if err != nil {
			return err
		}
		org, err := cmd.Flags().GetString("org")
		if err != nil {
			return err
		}

		if !withOrgs && org != "" {
			return errors.New("you can't use the --org flag without --with-orgs")
		}

		if withOrgs && org == "" {
			cmd.Println("You didn't specify an organization, fetching pull requests from all organizations")
			org = "all"
		}

		pullRequests, err := internal.FetchPullRequests(status, gitToken, repo, withOrgs, org)
		if err != nil {
			return err
		}

		internal.PrintPullRequests(pullRequests)

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func validateCommandArgs() error {
	err := isValidStatusArg(status)
	if err != nil {
		return err
	}

	if !withOrgs && org != "" {
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
	githubCmdFlags := rootCmd.Flags()
	githubCmdFlags.StringP("env-var", "", "GITHUB_TOKEN", "The name of the environment variable to fetch github's token")
	githubCmdFlags.StringP("repo", "r", "all", "The name of the repository to fetch pull requests from [all|<repo_name>]")
	githubCmdFlags.StringVarP(&status, "status", "s", "open", "The status of the pull requests to be fetched [open|closed|all]")
	githubCmdFlags.BoolVarP(&withOrgs, "with-orgs", "", false, "Fetch pull requests from organizations the user is part of")
	githubCmdFlags.StringVarP(&org, "org", "o", "", "The name of the organizations to fetch pull requests from [all|<org_name>]")
}
