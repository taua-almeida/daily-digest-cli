package internal

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/v56/github"
)

type repoError struct {
	repoName string
	err      error
}

func FetchPullRequests(args *Args, cmdConfig Config) ([]DetailedPullRequest, error) {
	token, err := GetGitToken(args.EnvVar)
	if err != nil {
		return nil, err
	}

	client := createClient(token)
	ctx := context.Background()
	user, err := getUser(ctx, client, cmdConfig.RateConfig)
	if err != nil {
		return nil, err
	}

	reposCollection, err := getReposCollection(ctx, client, args, user)

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	pullRequestsChan := make(chan []*github.PullRequest)
	errChan := make(chan repoError)

	for _, repoItem := range reposCollection {
		owner := repoItem.Owner
		for _, repName := range repoItem.Repositories {
			wg.Add(1)
			go func(repName string, owner string) {
				defer wg.Done()
				repoPr, _, err := client.PullRequests.List(ctx, owner, repName, &github.PullRequestListOptions{State: args.Status})
				if err != nil {
					fmt.Printf("Error fetching pull requests for repo %s: %v\n", repName, err)
					errChan <- repoError{repoName: repName, err: err}
					return
				}
				pullRequestsChan <- repoPr
			}(repName, owner)

		}

	}

	go func() {
		wg.Wait()
		close(pullRequestsChan)
		close(errChan)
	}()

	var allPullRequests []*github.PullRequest
	for pr := range pullRequestsChan {
		allPullRequests = append(allPullRequests, pr...)
	}

	var hasError bool
	for repoErr := range errChan {
		fmt.Printf("Error in repository %s: %v\n", repoErr.repoName, repoErr.err)
		hasError = true
	}

	if hasError {
		return nil, fmt.Errorf("errors occurred while fetching pull requests")
	}

	var detailedPrs []DetailedPullRequest

	for _, pr := range allPullRequests {
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
			detailedPrs = append(detailedPrs, detailedPR)
		}

	}

	return detailedPrs, nil
}
