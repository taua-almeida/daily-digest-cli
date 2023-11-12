package internal

import (
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
