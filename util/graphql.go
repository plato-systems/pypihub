package util

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func NewGitHubv4Client(ctx context.Context, token string) *githubv4.Client {
	var c *http.Client
	if TestGitHubAPI == nil {
		c = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
	} else {
		c = &http.Client{Transport: &testTransport}
	}
	return githubv4.NewClient(c)
}

var TestGitHubAPI http.HandlerFunc

type testTripper struct{}

func (t *testTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	TestGitHubAPI(w, req)
	return w.Result(), nil
}

var testTransport testTripper
