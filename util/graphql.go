package util

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// NewGitHubv4Client constructs a GraphQL client with the appropriate backend.
func NewGitHubv4Client(ctx context.Context, token string) *githubv4.Client {
	return githubv4.NewClient(oauth2.NewClient(
		ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	))
}
