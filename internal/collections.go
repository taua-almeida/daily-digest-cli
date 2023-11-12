package internal

import (
	"context"
	"fmt"

	"github.com/google/go-github/v56/github"
)

func getReposCollection(ctx context.Context, client *github.Client, args *Args, user *github.User) ([]RepositoryCollection, error) {
	var reposCollection []RepositoryCollection

	if args.WithOrgs {
		orgCollections, err := getOrgRepos(ctx, client, args.Org, user)
		if err != nil {
			return nil, err
		}
		reposCollection = append(reposCollection, orgCollections...)
	}

	personalRepos, err := getPersonalRepos(ctx, client, args.RepoName, user)

	if err != nil {
		return nil, err
	}

	if len(personalRepos) > 0 {
		reposCollection = append(reposCollection, RepositoryCollection{Repositories: personalRepos, Owner: *user.Login})
	}

	return reposCollection, nil
}

func getOrgRepos(ctx context.Context, client *github.Client, orgName string, user *github.User) ([]RepositoryCollection, error) {
	var orgCollections []RepositoryCollection

	if orgName == "all" {
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
			orgCollections = append(orgCollections, RepositoryCollection{Repositories: names, Owner: orgLogin})
		}

	} else {
		orgRepos, resp, err := client.Repositories.ListByOrg(ctx, orgName, &github.RepositoryListByOrgOptions{Type: "all"})
		if resp != nil && resp.StatusCode == 404 {
			return nil, fmt.Errorf("organization %s not found", orgName)
		}
		if err != nil {
			return nil, err
		}
		names := make([]string, 0, len(orgRepos))
		for _, repo := range orgRepos {
			names = append(names, *repo.Name)
		}
		orgCollections = append(orgCollections, RepositoryCollection{Repositories: names, Owner: orgName})
	}

	return orgCollections, nil
}

func getPersonalRepos(ctx context.Context, client *github.Client, repoName string, user *github.User) ([]string, error) {
	var personalRepos []string

	if repoName == "all" {
		repos, _, err := client.Repositories.List(ctx, *user.Login, &github.RepositoryListOptions{Affiliation: "owner,collaborator"})

		if err != nil {
			return nil, err
		}

		names := make([]string, 0, len(repos))
		for _, repo := range repos {
			names = append(names, *repo.Name)
		}
		personalRepos = append(personalRepos, names...)
	} else {
		repo, _, err := client.Repositories.Get(ctx, *user.Login, repoName)
		if err != nil {
			return nil, err
		}
		personalRepos = append(personalRepos, *repo.Name)
	}

	return personalRepos, nil
}
