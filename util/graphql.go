package util

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// GHv4Client represents a GitHub GraphQL API client.
type GHv4Client interface {
	Query(
		ctx context.Context, q interface{}, v map[string]interface{},
	) error
}

// GHv4ClientMaker describes how to create a GHv4Client.
type GHv4ClientMaker func(ctx context.Context, token string) GHv4Client

// NewGHv4Client constructs a production GitHub GraphQL API client.
func NewGHv4Client(ctx context.Context, token string) GHv4Client {
	return githubv4.NewClient(oauth2.NewClient(
		ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	))
}
