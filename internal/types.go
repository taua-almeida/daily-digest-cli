package internal

// Args holds the arguments passed to the application.
type Args struct {
	Status   string
	WithOrgs bool
	Org      string
	EnvVar   string
	RepoName string
}

// RepositoryCollection holds a collection of repositories for a user or organization.
type RepositoryCollection struct {
	Repositories []string // Just the names of the repositories
	Owner        string   // Owner login of the repositories (user or organization)
}

// DetailedPullRequest contains specific details of a pull request required by the application.
type DetailedPullRequest struct {
	Number      int    // Pull request number
	Title       string // Title of the pull request
	URL         string // HTML URL of the pull request
	State       string // State of the pull request (open, closed)
	Author      string // Author of the pull request
	IsMergeable bool   // Indicates if the pull request is mergeable
	CICDStatus  string // Continuous Integration/Continuous Deployment status
	Condition   string // Role of the authenticated user (author, assignee, reviewer)
}
