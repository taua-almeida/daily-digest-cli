package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v56/github"
)

func GetGitToken(envName string) (string, error) {
	token := os.Getenv(envName)

	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN environment variable not set or wrong flag used")
	}
	return token, nil
}

func createClient(token string) *github.Client {
	return github.NewClient(nil).WithAuthToken(token)
}

func checkRepoExists(ctx context.Context, client *github.Client, owner string, repo string) error {
	_, _, err := client.Repositories.Get(ctx, owner, repo)
	return err
}
