package util

import (
	"context"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// APIClient represents a GraphQL client with per-call authentication.
type APIClient interface {
	Query(
		ctx context.Context, token string,
		q interface{}, v map[string]interface{},
	) error
}

// GHv4Client is an APIClient exposing GitHub's GraphQL API.
type GHv4Client struct{}

func (g GHv4Client) Query(
	ctx context.Context, token string,
	q interface{}, v map[string]interface{},
) error {
	return githubv4.NewClient(oauth2.NewClient(
		ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	)).Query(ctx, q, v)
}
