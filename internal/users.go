package internal

import (
	"context"
	"fmt"

	"github.com/google/go-github/v56/github"
)

// getUser fetches a GitHub user using the provided client and token. It also checks
// if the request exceeds the rate limit set in rateConfig.
func getUser(ctx context.Context, client *github.Client, rateConfig RateConfig) (*github.User, error) {
	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	rateLimitExceeded := isRateLimitExceeded(resp, rateConfig)
	if rateLimitExceeded {
		message := fmt.Sprintf("Rate limit exceeded: %d requests remaining out of %d", resp.Rate.Remaining, resp.Rate.Limit)
		return nil, fmt.Errorf(message)
	}

	fmt.Printf("Rate remaining: %d\n", resp.Rate.Remaining)
	return user, nil
}

// isRateLimitExceeded checks if the remaining requests are below the configured rate limit.
func isRateLimitExceeded(resp *github.Response, rateConfig RateConfig) bool {
	var threshold int
	if rateConfig.RateType == "percentage" {
		threshold = rateConfig.RateLimit * resp.Rate.Limit / 100
	} else {
		threshold = rateConfig.RateLimit
	}

	return resp.Rate.Remaining < threshold
}

// getAllUserOrgsLogin fetches all active organization logins for the specified user.
func getAllUserOrgsLogin(ctx context.Context, client *github.Client, user string) ([]string, error) {
	orgs, _, err := client.Organizations.ListOrgMemberships(ctx, &github.ListOrgMembershipsOptions{State: "active"})
	if err != nil {
		return nil, err
	}

	orgLogins := make([]string, len(orgs))
	for i, org := range orgs {
		orgLogins[i] = *org.Organization.Login
	}

	return orgLogins, nil
}
