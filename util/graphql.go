package util

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// NewGitHubv4Client creates a new GitHub GraphQL API client with the given
// token as authorization.
func NewGitHubv4Client(ctx context.Context, token string) *githubv4.Client {
	hc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	return githubv4.NewClient(hc)
}
