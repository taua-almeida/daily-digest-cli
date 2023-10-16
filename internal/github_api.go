package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v56/github"
	"github.com/jedib0t/go-pretty/v6/table"
)

func GetGitToken(envName string) string {
	token := os.Getenv(envName)

	if token == "" {
		fmt.Println("Please set the GITHUB_TOKEN environment variable or use the right env flag")
		os.Exit(1)
	}
	return token
}

func FetchPullRequests(status string, token string, repo string, withOrg bool, org string) ([]*github.PullRequest, error) {
	client := createClient(token)
	ctx := context.Background()
	user, err := getUser(ctx, client)

	if err != nil {
		return nil, err
	}

	var repos []*github.Repository
	var pullRequests []*github.PullRequest

	if repo != "all" {

		// Checking if repo exists
		err = checkRepoExists(ctx, client, *user.Login, repo)
		if err != nil {
			return nil, err
		}

		// Fetching pull requests from a single repo
		pullRequests, _, err = client.PullRequests.List(ctx, *user.Login, repo, &github.PullRequestListOptions{State: status})
		if err != nil {
			return nil, err
		}

		// TODO with orgs and repo name should be possible

		return pullRequests, nil

	}

	if !withOrg {
		repos, _, err = client.Repositories.List(ctx, *user.Login, &github.RepositoryListOptions{Affiliation: "owner,collaborator"})
		if err != nil {
			return nil, err
		}
	} else {
		if org == "all" {
			orgsLogin, err := getAllUserOrgsLogin(ctx, client, *user.Login)
			if err != nil {
				return nil, err
			}
			for _, orgLogin := range orgsLogin {
				orgRepos, _, err := client.Repositories.ListByOrg(ctx, orgLogin, &github.RepositoryListByOrgOptions{Type: "all"})
				if err != nil {
					return nil, err
				}

				repos = append(repos, orgRepos...)
			}

		} else {
			orgRepos, _, err := client.Repositories.ListByOrg(ctx, org, &github.RepositoryListByOrgOptions{Type: "all"})
			if err != nil {
				return nil, err
			}
			repos = append(repos, orgRepos...)
		}
		personalRepo, _, err := client.Repositories.List(ctx, *user.Login, &github.RepositoryListOptions{})

		if err != nil {
			return nil, err
		}

		repos = append(repos, personalRepo...)
	}

	for _, repo := range repos {
		repoPr, _, err := client.PullRequests.List(ctx, *user.Login, *repo.Name, &github.PullRequestListOptions{State: status})
		if err != nil {
			return nil, err
		}
		pullRequests = append(pullRequests, repoPr...)
	}

	return pullRequests, nil
}

func PrintPullRequests(pullRequests []*github.PullRequest) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"PR Number", "Title", "URL", "Status", "MergeableState"})
	for _, pr := range pullRequests {
		if pr == nil {
			continue // Skip this iteration if pr is nil
		}

		// Safely get the values, using defaults or placeholders if nil
		number := 0
		if pr.Number != nil {
			number = *pr.Number
		}

		title := "N/A"
		if pr.Title != nil {
			title = *pr.Title
		}

		url := "N/A"
		if pr.URL != nil {
			url = *pr.URL
		}

		state := "N/A"
		if pr.State != nil {
			state = *pr.State
		}

		mergeableState := "N/A"
		if pr.MergeableState != nil {
			mergeableState = *pr.MergeableState
		}

		// Append the row with the safely extracted values
		t.AppendRow(table.Row{number, title, url, state, mergeableState})
	}

	t.AppendFooter(table.Row{"Total", len(pullRequests), "", "", ""})
	t.SetStyle(table.StyleColoredBright)
	t.Render()
}

// Gets all user organizations login names by using ListOrgMemberships
func getAllUserOrgsLogin(ctx context.Context, client *github.Client, user string) ([]string, error) {
	orgs, _, err := client.Organizations.ListOrgMemberships(ctx, &github.ListOrgMembershipsOptions{State: "active"})
	if err != nil {
		return nil, err
	}
	orgsLogin := make([]string, 0, len(orgs))
	for _, org := range orgs {
		orgsLogin = append(orgsLogin, *org.Organization.Login)
	}
	return orgsLogin, nil
}

// Checks if a repo name exists
func checkRepoExists(ctx context.Context, client *github.Client, owner string, repo string) error {
	_, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return err
	}
	return nil
}

// Gets the user that the token belongs to
// If the rate limit is 10% or less, it stops the execution
func getUser(ctx context.Context, client *github.Client) (*github.User, error) {
	rateLimitPercentage := 0.1

	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	if resp.Rate.Remaining < int(rateLimitPercentage*float64(resp.Rate.Limit)) {
		fmt.Printf("You have %d requests remaining out of %d, stopping execution", resp.Rate.Remaining, resp.Rate.Limit)
		return nil, fmt.Errorf("rate limit is 10%% or less, stopping execution")
	}

	return user, nil
}

func createClient(token string) *github.Client {
	return github.NewClient(nil).WithAuthToken(token)
}
