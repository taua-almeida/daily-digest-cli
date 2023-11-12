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

	var pullRequests []*github.PullRequest
	var filteredPullRequests []DetailedPullRequest

	reposCollection, err := getReposCollection(ctx, client, args, user)

	if err != nil {
		return nil, err
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
