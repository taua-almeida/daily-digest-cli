package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v56/github"
	"github.com/jedib0t/go-pretty/v6/table"
)

type RepoStruct struct {
	repos []*github.Repository
	login string
}

type PullRequestStruct struct {
	pr         *github.PullRequest
	condition  string
	mergeState bool
	ciCDStatus string
}

func GetGitToken(envName string) string {
	token := os.Getenv(envName)

	if token == "" {
		fmt.Println("Please set the GITHUB_TOKEN environment variable or use the right env flag")
		os.Exit(1)
	}
	return token
}

func FetchPullRequests(status string, token string, repo string, withOrg bool, org string) ([]PullRequestStruct, error) {
	client := createClient(token)
	ctx := context.Background()
	user, err := getUser(ctx, client)

	if err != nil {
		return nil, err
	}

	var repos []RepoStruct
	var pullRequests []*github.PullRequest
	var filteredPullRequests []PullRequestStruct

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

		for _, pr := range pullRequests {
			if pr.User != nil && pr.User.Login != nil && *pr.User.Login == *user.Login {
				filteredPullRequests = append(filteredPullRequests, PullRequestStruct{pr: pr, condition: "author", mergeState: false})
			}

			if pr.Assignee != nil && pr.Assignee.Login != nil && *pr.Assignee.Login == *user.Login {
				filteredPullRequests = append(filteredPullRequests, PullRequestStruct{pr: pr, condition: "assignee", mergeState: false})
			}

			if pr.RequestedReviewers != nil {
				for _, reviewer := range pr.RequestedReviewers {
					if reviewer.Login != nil && *reviewer.Login == *user.Login {
						filteredPullRequests = append(filteredPullRequests, PullRequestStruct{pr: pr, condition: "reviewer", mergeState: false})
					}
				}
			}
		}
		// TODO: remove this return, because we can have a single repo with the org condition too
		return filteredPullRequests, nil

	}

	if !withOrg {
		clientRepos, _, err := client.Repositories.List(ctx, *user.Login, &github.RepositoryListOptions{Affiliation: "owner,collaborator"})
		if err != nil {
			return nil, err
		}
		repos = append(repos, RepoStruct{repos: clientRepos, login: *user.Login})
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

				repos = append(repos, RepoStruct{repos: orgRepos, login: orgLogin})
			}

		} else {
			orgRepos, _, err := client.Repositories.ListByOrg(ctx, org, &github.RepositoryListByOrgOptions{Type: "all"})
			if err != nil {
				return nil, err
			}
			repos = append(repos, RepoStruct{repos: orgRepos, login: org})
		}
		personalRepo, _, err := client.Repositories.List(ctx, *user.Login, &github.RepositoryListOptions{})

		if err != nil {
			return nil, err
		}

		repos = append(repos, RepoStruct{repos: personalRepo, login: *user.Login})
	}

	for _, repoItem := range repos {
		for _, rep := range repoItem.repos {
			repoPr, _, err := client.PullRequests.List(ctx, repoItem.login, *rep.Name, &github.PullRequestListOptions{State: status})
			if err != nil {
				return nil, err
			}
			pullRequests = append(pullRequests, repoPr...)
		}

	}

	for _, pr := range pullRequests {
		if pr.User != nil && pr.User.Login != nil && *pr.User.Login == *user.Login {
			commitStatus, _, err := client.Repositories.GetCombinedStatus(ctx, *pr.Base.Repo.Owner.Login, *pr.Base.Repo.Name, *pr.Head.SHA, nil)
			if err != nil {
				return nil, err
			}
			filteredPullRequests = append(filteredPullRequests, PullRequestStruct{pr: pr, condition: "author", mergeState: pr.GetMergeable(), ciCDStatus: commitStatus.GetState()})
		}

		if pr.Assignee != nil && pr.Assignee.Login != nil && *pr.Assignee.Login == *user.Login {
			commitStatus, _, err := client.Repositories.GetCombinedStatus(ctx, *pr.Base.Repo.Owner.Login, *pr.Base.Repo.Name, *pr.Head.SHA, nil)
			if err != nil {
				return nil, err
			}
			filteredPullRequests = append(filteredPullRequests, PullRequestStruct{pr: pr, condition: "assignee", mergeState: pr.GetMergeable(), ciCDStatus: commitStatus.GetState()})
		}

		if pr.RequestedReviewers != nil {
			for _, reviewer := range pr.RequestedReviewers {
				if reviewer.Login != nil && *reviewer.Login == *user.Login {
					commitStatus, _, err := client.Repositories.GetCombinedStatus(ctx, *pr.Base.Repo.Owner.Login, *pr.Base.Repo.Name, *pr.Head.SHA, nil)
					if err != nil {
						return nil, err
					}
					filteredPullRequests = append(filteredPullRequests, PullRequestStruct{pr: pr, condition: "reviewer", mergeState: pr.GetMergeable(), ciCDStatus: commitStatus.GetState()})
				}
			}
		}

	}

	return filteredPullRequests, nil
}

func PrintPullRequests(pullRequests []PullRequestStruct) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"PR Number", "Title", "PR Link", "Status", "PR Condition", "Mergeable", "CI/CD Status"})
	for _, prItem := range pullRequests {
		// Safely get the values, using defaults or placeholders if nil
		number := 0
		if prItem.pr.Number != nil {
			number = *prItem.pr.Number
		}

		title := "N/A"
		if prItem.pr.Title != nil {
			title = *prItem.pr.Title
		}

		prLink := "N/A"
		if prItem.pr.HTMLURL != nil {
			prLink = *prItem.pr.HTMLURL
		}

		state := "N/A"
		if prItem.pr.State != nil {
			state = *prItem.pr.State
		}

		// Append the row with the safely extracted values
		t.AppendRow(table.Row{number, title, prLink, state, prItem.condition, prItem.mergeState, prItem.ciCDStatus})
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

	fmt.Println("Rate remaining:", resp.Rate.Remaining)

	if resp.Rate.Remaining < int(rateLimitPercentage*float64(resp.Rate.Limit)) {
		fmt.Printf("You have %d requests remaining out of %d, stopping execution", resp.Rate.Remaining, resp.Rate.Limit)
		return nil, fmt.Errorf("rate limit is 10%% or less, stopping execution")
	}

	return user, nil
}

func createClient(token string) *github.Client {
	return github.NewClient(nil).WithAuthToken(token)
}
