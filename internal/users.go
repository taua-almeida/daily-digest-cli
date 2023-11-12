package internal

import (
	"context"
	"fmt"

	"github.com/google/go-github/v56/github"
)

// Get a user from the github Token and reports rate limit
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
