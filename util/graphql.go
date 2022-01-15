package util

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GHv4ClientMaker func(ctx context.Context, token string) *githubv4.Client

// NewGHv4Client constructs a GraphQL client for GitHub.
// TODO: encapsulate token in ctx?
func NewGHv4Client(ctx context.Context, token string) *githubv4.Client {
	return githubv4.NewClient(oauth2.NewClient(
		ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	))
}

type testTripper struct {
	serve http.HandlerFunc
}

func (t testTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	t.serve(w, req)
	return w.Result(), nil
}

func NewGHv4ClientMaker(serve http.HandlerFunc) GHv4ClientMaker {
	c := githubv4.NewClient(&http.Client{Transport: testTripper{serve}})
	return func(context.Context, string) *githubv4.Client {
		return c
	}
}
