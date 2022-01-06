package util

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type APIClient interface {
	Query(
		ctx context.Context, q interface{},
		variables map[string]interface{},
	) error
}

type APIClientFactory interface {
	New(ctx context.Context, token string) APIClient
}

type GitHubv4ClientFactory struct{}

func (c GitHubv4ClientFactory) New(
	ctx context.Context, token string,
) APIClient {
	hc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	return githubv4.NewClient(hc)
}
