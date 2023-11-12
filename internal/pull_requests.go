package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v56/github"
)

func FetchPullRequests(args *Args) ([]DetailedPullRequest, error) {
	start := time.Now()
	token, err := GetGitToken(args.EnvVar)
	if err != nil {
		return nil, err
	}
	client := createClient(token)
	ctx := context.Background()
	user, err := getUser(ctx, client)

	if err != nil {
		return nil, err
	}

	var reposCollection []RepositoryCollection
	var pullRequests []*github.PullRequest
	var filteredPullRequests []DetailedPullRequest

	if !args.WithOrgs {
		if args.RepoName != "all" {
			// Checking if repo exists
			err = checkRepoExists(ctx, client, *user.Login, args.RepoName)
			if err != nil {
				return nil, err
			}

			// Fetching pull requests from a single repo
			pullRequests, _, err = client.PullRequests.List(ctx, *user.Login, args.RepoName, &github.PullRequestListOptions{State: args.Status})
			if err != nil {
				return nil, err
			}
			pullRequests = append(pullRequests, pullRequests...)

		} else {
			clientRepos, _, err := client.Repositories.List(ctx, *user.Login, &github.RepositoryListOptions{Affiliation: "owner,collaborator"})
			if err != nil {
				return nil, err
			}
			names := make([]string, 0, len(clientRepos))
			for _, repo := range clientRepos {
				names = append(names, *repo.Name)
			}
			reposCollection = append(reposCollection, RepositoryCollection{Repositories: names, Owner: *user.Login})
		}

	} else {
		if args.Org == "all" {
			orgsLogin, err := getAllUserOrgsLogin(ctx, client, *user.Login)
			if err != nil {
				return nil, err
			}
			for _, orgLogin := range orgsLogin {
				orgRepos, _, err := client.Repositories.ListByOrg(ctx, orgLogin, &github.RepositoryListByOrgOptions{Type: "all"})
				if err != nil {
					return nil, err
				}

				names := make([]string, 0, len(orgRepos))
				for _, repo := range orgRepos {
					names = append(names, repo.GetName())
				}
				reposCollection = append(reposCollection, RepositoryCollection{Repositories: names, Owner: orgLogin})
			}

		} else {
			orgRepos, resp, err := client.Repositories.ListByOrg(ctx, args.Org, &github.RepositoryListByOrgOptions{Type: "all"})
			if resp != nil && resp.StatusCode == 404 {
				return nil, fmt.Errorf("organization %s not found", args.Org)
			}
			if err != nil {
				return nil, err
			}
			names := make([]string, 0, len(orgRepos))
			for _, repo := range orgRepos {
				names = append(names, *repo.Name)
			}
			reposCollection = append(reposCollection, RepositoryCollection{Repositories: names, Owner: args.Org})
		}

		if args.RepoName != "all" {
			// Checking if repo exists
			err = checkRepoExists(ctx, client, *user.Login, args.RepoName)
			if err != nil {
				return nil, err
			}

			// Fetching pull requests from a single repo
			pullRequests, _, err = client.PullRequests.List(ctx, *user.Login, args.RepoName, &github.PullRequestListOptions{State: args.Status})
			if err != nil {
				return nil, err
			}
			pullRequests = append(pullRequests, pullRequests...)

		} else {
			personalRepo, _, err := client.Repositories.List(ctx, *user.Login, &github.RepositoryListOptions{Affiliation: "owner,collaborator"})

			if err != nil {
				return nil, err
			}

			names := make([]string, 0, len(personalRepo))
			for _, repo := range personalRepo {
				names = append(names, *repo.Name)
			}
			reposCollection = append(reposCollection, RepositoryCollection{Repositories: names, Owner: *user.Login})
		}

	}

	for _, repoItem := range reposCollection {
		for _, repName := range repoItem.Repositories {
			repoPr, _, err := client.PullRequests.List(ctx, repoItem.Owner, repName, &github.PullRequestListOptions{State: args.Status})
			if err != nil {
				return nil, err
			}
			pullRequests = append(pullRequests, repoPr...)
		}

	}

	for _, pr := range pullRequests {
		detailedPR := DetailedPullRequest{
			Number:      pr.GetNumber(),
			Title:       pr.GetTitle(),
			URL:         pr.GetHTMLURL(),
			State:       pr.GetState(),
			Author:      pr.User.GetLogin(),
			IsMergeable: pr.GetMergeable(),
			CICDStatus:  "",
			Condition:   "",
		}

		if pr.User != nil && pr.User.Login != nil && *pr.User.Login == *user.Login {
			commitStatus, _, err := client.Repositories.GetCombinedStatus(ctx, *pr.Base.Repo.Owner.Login, *pr.Base.Repo.Name, *pr.Head.SHA, nil)
			if err != nil {
				return nil, err
			}
			detailedPR.Condition = "author"
			detailedPR.CICDStatus = commitStatus.GetState()
		}

		if pr.Assignee != nil && pr.Assignee.Login != nil && *pr.Assignee.Login == *user.Login {
			commitStatus, _, err := client.Repositories.GetCombinedStatus(ctx, *pr.Base.Repo.Owner.Login, *pr.Base.Repo.Name, *pr.Head.SHA, nil)
			if err != nil {
				return nil, err
			}
			detailedPR.Condition = "assignee"
			detailedPR.CICDStatus = commitStatus.GetState()
		}

		if pr.RequestedReviewers != nil {
			for _, reviewer := range pr.RequestedReviewers {
				if reviewer.Login != nil && *reviewer.Login == *user.Login {
					commitStatus, _, err := client.Repositories.GetCombinedStatus(ctx, *pr.Base.Repo.Owner.Login, *pr.Base.Repo.Name, *pr.Head.SHA, nil)
					if err != nil {
						return nil, err
					}
					detailedPR.Condition = "reviewer"
					detailedPR.CICDStatus = commitStatus.GetState()
				}
			}
		}

		if detailedPR.Condition != "" {
			filteredPullRequests = append(filteredPullRequests, detailedPR)
		}

	}

	elapsed := time.Since(start)
	fmt.Printf("Code block took %s\n", elapsed)

	return filteredPullRequests, nil
}
